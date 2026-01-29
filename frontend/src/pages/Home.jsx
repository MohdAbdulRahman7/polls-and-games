import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import { toast } from 'react-toastify'

function Home() {
  const [polls, setPolls] = useState([])
  const [page, setPage] = useState(1)
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const limit = 10

  useEffect(() => {
    fetchPolls()
  }, [page])

  const fetchPolls = async () => {
    try {
      setLoading(true)
      const response = await axios.get(`/api/polls?page=${page}`)
      setPolls(response.data.polls || [])
      setTotal(response.data.total || 0)
    } catch (error) {
      toast.error('Failed to fetch polls')
      setPolls([])
      setTotal(0)
    } finally {
      setLoading(false)
    }
  }

  const totalPages = Math.ceil(total / limit)

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="text-center mb-8">
        <h1 className="text-4xl font-bold mb-4">All Polls</h1>
        <p className="text-lg text-base-content/70">Browse and participate in polls created by the community</p>
      </div>

      {loading ? (
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      ) : !polls || polls.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-xl text-base-content/70">No polls available yet. Be the first to create one!</p>
        </div>
      ) : (
        <>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
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

          {totalPages > 1 && (
            <div className="flex justify-center gap-2">
              <button
                className="btn btn-outline"
                onClick={() => setPage((p) => Math.max(1, p - 1))}
                disabled={page === 1}
              >
                Previous
              </button>
              <span className="flex items-center px-4">
                Page {page} of {totalPages}
              </span>
              <button
                className="btn btn-outline"
                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                disabled={page === totalPages}
              >
                Next
              </button>
            </div>
          )}
        </>
      )}
    </div>
  )
}

export default Home

