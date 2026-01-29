import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { toast } from 'react-toastify'
import { useAuth } from '../context/AuthContext'

function Dashboard() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [polls, setPolls] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      navigate('/login')
      return
    }
    fetchUserPolls()
  }, [user])

  const fetchUserPolls = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`/api/user/${user.id}/polls`)
      setPolls(response.data)
    } catch (error) {
      toast.error('Failed to fetch your polls')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (pollId) => {
    if (!window.confirm('Are you sure you want to delete this poll?')) {
      return
    }

    try {
      await axios.delete(`/api/polls/${pollId}`)
      toast.success('Poll deleted successfully')
      fetchUserPolls()
    } catch (error) {
      toast.error('Failed to delete poll')
    }
  }

  if (!user) {
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="card bg-base-100 shadow-xl mb-8">
        <div className="card-body">
          <h1 className="card-title text-3xl mb-4">Dashboard</h1>
          <div className="space-y-2">
            <p><strong>Username:</strong> {user.username}</p>
            <p><strong>Email:</strong> {user.email}</p>
            <p><strong>User ID:</strong> {user.id}</p>
          </div>
        </div>
      </div>

      <div className="flex justify-between items-center mb-6">
        <h2 className="text-2xl font-bold">My Polls</h2>
        <Link to="/create-poll" className="btn btn-primary">
          Create New Poll
        </Link>
      </div>

      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : polls.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-xl text-base-content/70 mb-4">You haven't created any polls yet.</p>
          <Link to="/create-poll" className="btn btn-primary">
            Create Your First Poll
          </Link>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {polls.map((poll) => (
            <div key={poll.id} className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">{poll.title}</h2>
                <p className="text-sm text-base-content/70 line-clamp-2">{poll.description}</p>
                <div className="flex items-center justify-between mt-4">
                  <div className="text-sm text-base-content/60">
                    <span>{poll.vote_count} votes</span>
                    <span className="mx-2">•</span>
                    <span>{poll.options.length} options</span>
                  </div>
                </div>
                <div className="card-actions justify-end mt-4">
                  <Link to={`/poll/${poll.id}`} className="btn btn-primary btn-sm">
                    View
                  </Link>
                  <button
                    className="btn btn-error btn-sm"
                    onClick={() => handleDelete(poll.id)}
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default Dashboard

