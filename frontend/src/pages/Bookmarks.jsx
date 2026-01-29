import { useState, useEffect } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import axios from 'axios'
import { toast } from 'react-toastify'
import { useAuth } from '../context/AuthContext'

function Bookmarks() {
  const { user } = useAuth()
  const navigate = useNavigate()
  const [polls, setPolls] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      navigate('/login')
      return
    }
    fetchBookmarks()
  }, [user])

  const fetchBookmarks = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`/api/bookmarks/${user.id}`)
      setPolls(response.data)
    } catch (error) {
      toast.error('Failed to fetch bookmarks')
    } finally {
      setLoading(false)
    }
  }

  if (!user) {
    return null
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">Bookmarked Polls</h1>

      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : polls.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-xl text-base-content/70 mb-4">You haven't bookmarked any polls yet.</p>
          <Link to="/" className="btn btn-primary">
            Browse Polls
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
                    <span>By {poll.username}</span>
                    <span className="mx-2">•</span>
                    <span>{poll.vote_count} votes</span>
                  </div>
                  <Link to={`/poll/${poll.id}`} className="btn btn-primary btn-sm">
                    View Poll
                  </Link>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}

export default Bookmarks

