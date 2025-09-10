import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute('/')({
  component: App,
})

function App() {
  return (
    <div className="container mx-auto p-6">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-gray-900 mb-4">
            Welcome to BriefBot
          </h1>
          <p className="text-xl text-gray-600 mb-8">
            An AI-enabled platform for managing links and extracting knowledge faster
          </p>
        </div>

        <div className="grid md:grid-cols-2 gap-8 mb-12">
          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">User Management</h2>
            <p className="text-gray-600 mb-6">
              Create and manage users in the system. Users can store and organize their content items.
            </p>
            <Link
              to="/users"
              className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 transition-colors"
            >
              Manage Users
            </Link>
          </div>

          <div className="bg-white p-6 rounded-lg shadow-md">
            <h2 className="text-2xl font-semibold text-gray-900 mb-4">Item Management</h2>
            <p className="text-gray-600 mb-6">
              Create and manage content items. Items represent URLs with extracted text content that users can read and organize.
            </p>
            <Link
              to="/items"
              className="inline-flex items-center px-4 py-2 bg-green-600 text-white rounded-md hover:bg-green-700 transition-colors"
            >
              Manage Items
            </Link>
          </div>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-md">
          <h2 className="text-2xl font-semibold text-gray-900 mb-4">Quick Start</h2>
          <ol className="list-decimal list-inside space-y-2 text-gray-700">
            <li>Go to <Link to="/users" className="text-blue-600 hover:underline">Users</Link> and create a new user</li>
            <li>Go to <Link to="/items" className="text-blue-600 hover:underline">Items</Link> and create items for that user</li>
            <li>View and manage items with read/unread status</li>
            <li>Update or delete users and items as needed</li>
          </ol>
        </div>

        <div className="mt-8 text-center">
          <p className="text-gray-500">
            Built with React, Tanstack Router & Query, Go, Gin, and PostgreSQL
          </p>
        </div>
      </div>
    </div>
  )
}
