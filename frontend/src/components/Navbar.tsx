import { Link } from "@tanstack/react-router"
import { FileText, Podcast, Sparkles } from "lucide-react"

export function Navbar() {
  return (
    <nav className="border-b bg-white">
      <div className="container mx-auto px-4">
        <div className="flex h-16 items-center justify-between">
          {/* Logo/Brand */}
          <Link to="/" className="flex items-center gap-2 text-xl font-bold">
            <Sparkles className="h-6 w-6 text-primary" />
            <span>BriefBot</span>
          </Link>
          <a
            href="https://github.com/yamirghofran/briefbot"
            target="_blank"
            rel="noreferrer"
            className="ml-4 text-sm font-semibold text-gray-600 hover:text-gray-900"
          >
            View GitHub
          </a>

          {/* Navigation Links */}
          <div className="flex items-center gap-1">
            <NavLink to="/" icon={FileText}>
              Items
            </NavLink>
            <NavLink to="/items/podcasts" icon={Podcast}>
              Podcasts
            </NavLink>
          </div>
        </div>
      </div>
    </nav>
  )
}

interface NavLinkProps {
  to: string
  icon: React.ElementType
  children: React.ReactNode
}

function NavLink({ to, icon: Icon, children }: NavLinkProps) {
  return (
    <Link
      to={to}
      className="flex items-center gap-2 px-4 py-2 rounded-md text-sm font-medium transition-colors hover:bg-gray-100"
      activeProps={{
        className: "bg-gray-100 text-primary",
      }}
      inactiveProps={{
        className: "text-gray-600",
      }}
    >
      <Icon className="h-4 w-4" />
      {children}
    </Link>
  )
}
