import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import { toast } from 'react-toastify'
import { useAuth } from '../context/AuthContext'

function CreatePoll() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [options, setOptions] = useState(['', ''])
  const [submitting, setSubmitting] = useState(false)

  const addOption = () => {
    if (options.length < 10) {
      setOptions([...options, ''])
    } else {
      toast.error('Maximum 10 options allowed')
    }
  }

  const removeOption = (index) => {
    if (options.length > 2) {
      setOptions(options.filter((_, i) => i !== index))
    } else {
      toast.error('At least 2 options are required')
    }
  }

  const updateOption = (index, value) => {
    const newOptions = [...options]
    newOptions[index] = value
    setOptions(newOptions)
  }

  const handleSubmit = async (e) => {
    e.preventDefault()

    if (!title.trim()) {
      toast.error('Please enter a title')
      return
    }

    const validOptions = options.filter(opt => opt.trim())
    if (validOptions.length < 2) {
      toast.error('Please provide at least 2 options')
      return
    }

    try {
      setSubmitting(true)
      const response = await axios.post('/api/polls', {
        user_id: user.id,
        title: title.trim(),
        description: description.trim(),
        options: validOptions,
      })
      toast.success('Poll created successfully!')
      navigate(`/poll/${response.data.id}`)
    } catch (error) {
      toast.error('Failed to create poll')
    } finally {
      setSubmitting(false)
    }
  }

  if (!user) {
    navigate('/login')
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <h1 className="text-3xl font-bold mb-8">Create New Poll</h1>

      <form onSubmit={handleSubmit}>
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Poll Title *</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                placeholder="Enter poll title"
                required
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Description</span>
              </label>
              <textarea
                className="textarea textarea-bordered"
                value={description}
                onChange={(e) => setDescription(e.target.value)}
                placeholder="Enter poll description (optional)"
                rows="4"
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Options *</span>
              </label>
              <div className="space-y-2">
                {options.map((option, index) => (
                  <div key={index} className="flex gap-2">
                    <input
                      type="text"
                      className="input input-bordered flex-1"
                      value={option}
                      onChange={(e) => updateOption(index, e.target.value)}
                      placeholder={`Option ${index + 1}`}
                    />
                    {options.length > 2 && (
                      <button
                        type="button"
                        className="btn btn-error btn-sm"
                        onClick={() => removeOption(index)}
                      >
                        Remove
                      </button>
                    )}
                  </div>
                ))}
              </div>
              <button
                type="button"
                className="btn btn-outline btn-sm mt-2"
                onClick={addOption}
                disabled={options.length >= 10}
              >
                Add Option
              </button>
            </div>

            <div className="card-actions justify-end mt-6">
              <button
                type="button"
                className="btn btn-ghost"
                onClick={() => navigate('/')}
              >
                Cancel
              </button>
              <button
                type="submit"
                className="btn btn-primary"
                disabled={submitting}
              >
                {submitting ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Creating...
                  </>
                ) : (
                  'Create Poll'
                )}
              </button>
            </div>
          </div>
        </div>
      </form>
    </div>
  )
}

export default CreatePoll

