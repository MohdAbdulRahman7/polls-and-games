import { Link } from 'react-router-dom'
import { useAuth } from '../context/AuthContext'

function Navbar() {
  const { user, logout } = useAuth()

  return (
    <div className="navbar bg-base-100 shadow-lg">
      <div className="navbar-start">
        <Link to="/" className="btn btn-ghost normal-case text-xl">
          Polls & Games
        </Link>
      </div>
      <div className="navbar-end">
        {user ? (
          <>
            <Link to="/" className="btn btn-ghost">
              Home
            </Link>
            <Link to="/create-poll" className="btn btn-primary">
              Create Poll
            </Link>
            <Link to="/bookmarks" className="btn btn-ghost">
              Bookmarks
            </Link>
            <Link to="/dashboard" className="btn btn-ghost">
              Dashboard
            </Link>
            <div className="dropdown dropdown-end">
              <label tabIndex={0} className="btn btn-ghost btn-circle avatar">
                <div className="w-10 rounded-full bg-primary text-primary-content flex items-center justify-center">
                  {user.username.charAt(0).toUpperCase()}
                </div>
              </label>
              <ul tabIndex={0} className="menu menu-sm dropdown-content mt-3 z-[1] p-2 shadow bg-base-100 rounded-box w-52">
                <li>
                  <span className="justify-between">
                    {user.username}
                    <span className="badge">User</span>
                  </span>
                </li>
                <li>
                  <a onClick={logout}>Logout</a>
                </li>
              </ul>
            </div>
          </>
        ) : (
          <>
            <Link to="/login" className="btn btn-ghost">
              Login
            </Link>
            <Link to="/register" className="btn btn-primary">
              Sign Up
            </Link>
          </>
        )}
      </div>
    </div>
  )
}

export default Navbar

