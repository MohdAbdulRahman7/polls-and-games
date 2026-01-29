import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { toast } from 'react-toastify'
import { useAuth } from '../context/AuthContext'

function PollView() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useAuth()
  const [poll, setPoll] = useState(null)
  const [selectedOption, setSelectedOption] = useState(null)
  const [bookmarked, setBookmarked] = useState(false)
  const [loading, setLoading] = useState(true)
  const [voting, setVoting] = useState(false)

  useEffect(() => {
    fetchPoll()
    if (user) {
      checkBookmark()
    }
  }, [id, user])

  const fetchPoll = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`/api/polls/${id}`)
      setPoll(response.data)
    } catch (error) {
      toast.error('Failed to fetch poll')
      navigate('/')
    } finally {
      setLoading(false)
    }
  }

  const checkBookmark = async () => {
    try {
      const response = await axios.get(`/api/check-bookmark?user_id=${user.id}&poll_id=${id}`)
      setBookmarked(response.data.bookmarked)
    } catch (error) {
      // Ignore error
    }
  }

  const handleVote = async () => {
    if (!user) {
      toast.error('Please login to vote')
      navigate('/login')
      return
    }

    if (!selectedOption) {
      toast.error('Please select an option')
      return
    }

    try {
      setVoting(true)
      await axios.post('/api/vote', {
        user_id: user.id,
        poll_id: parseInt(id),
        option_id: selectedOption,
      })
      toast.success('Vote recorded!')
      fetchPoll()
    } catch (error) {
      toast.error('Failed to vote')
    } finally {
      setVoting(false)
    }
  }

  const handleBookmark = async () => {
    if (!user) {
      toast.error('Please login to bookmark')
      navigate('/login')
      return
    }

    try {
      if (bookmarked) {
        await axios.delete('/api/bookmark', {
          data: {
            user_id: user.id,
            poll_id: parseInt(id),
          },
        })
        setBookmarked(false)
        toast.success('Removed from bookmarks')
      } else {
        await axios.post('/api/bookmark', {
          user_id: user.id,
          poll_id: parseInt(id),
        })
        setBookmarked(true)
        toast.success('Added to bookmarks')
      }
    } catch (error) {
      toast.error('Failed to update bookmark')
    }
  }

  const totalVotes = poll?.options.reduce((sum, opt) => sum + opt.vote_count, 0) || 0

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-[400px]">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    )
  }

  if (!poll) {
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-3xl">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex justify-between items-start mb-4">
            <div>
              <h1 className="card-title text-3xl mb-2">{poll.title}</h1>
              <p className="text-base-content/70 mb-4">{poll.description}</p>
              <div className="text-sm text-base-content/60">
                <span>Created by {poll.username}</span>
                <span className="mx-2">•</span>
                <span>{new Date(poll.created_at).toLocaleDateString()}</span>
                <span className="mx-2">•</span>
                <span>{totalVotes} total votes</span>
              </div>
            </div>
            {user && (
              <button
                className={`btn btn-circle ${bookmarked ? 'btn-primary' : 'btn-outline'}`}
                onClick={handleBookmark}
              >
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-6 w-6"
                  fill={bookmarked ? 'currentColor' : 'none'}
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 5a2 2 0 012-2h10a2 2 0 012 2v16l-7-3.5L5 21V5z"
                  />
                </svg>
              </button>
            )}
          </div>

          <div className="divider"></div>

          <div className="space-y-4">
            {poll.options.map((option) => {
              const percentage = totalVotes > 0 ? (option.vote_count / totalVotes) * 100 : 0
              return (
                <div key={option.id} className="form-control">
                  <label className="label cursor-pointer">
                    <div className="flex-1">
                      <div className="flex items-center justify-between mb-2">
                        <span className="label-text text-lg">{option.text}</span>
                        <span className="text-sm text-base-content/60">
                          {option.vote_count} votes ({percentage.toFixed(1)}%)
                        </span>
                      </div>
                      <progress
                        className="progress progress-primary w-full"
                        value={option.vote_count}
                        max={totalVotes || 1}
                      ></progress>
                    </div>
                    <input
                      type="radio"
                      name="poll-option"
                      className="radio radio-primary ml-4"
                      checked={selectedOption === option.id}
                      onChange={() => setSelectedOption(option.id)}
                    />
                  </label>
                </div>
              )
            })}
          </div>

          {user && (
            <div className="card-actions justify-end mt-6">
              <button
                className="btn btn-primary"
                onClick={handleVote}
                disabled={!selectedOption || voting}
              >
                {voting ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Voting...
                  </>
                ) : (
                  'Vote'
                )}
              </button>
            </div>
          )}

          {!user && (
            <div className="alert alert-info mt-6">
              <span>Please login to vote on this poll</span>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default PollView

