import { Link, Outlet } from '@tanstack/react-router'

export function Layout() {
  return (
    <div className="min-h-screen bg-gray-50">
      <nav className="bg-white shadow-sm border-b">
        <div className="container mx-auto px-6 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-8">
              <Link to="/" className="text-xl font-bold text-gray-900 hover:text-gray-700">
                BriefBot
              </Link>
              <div className="flex space-x-6">
                <Link
                  to="/users"
                  className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
                  activeProps={{
                    className: 'bg-gray-100 text-gray-900',
                  }}
                >
                  Users
                </Link>
                <Link
                  to="/items"
                  className="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium"
                  activeProps={{
                    className: 'bg-gray-100 text-gray-900',
                  }}
                >
                  Items
                </Link>
              </div>
            </div>
          </div>
        </div>
      </nav>
      <main>
        <Outlet />
      </main>
    </div>
  )
}