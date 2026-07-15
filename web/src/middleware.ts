import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

const publicRoutes = ['/', '/login', '/register', '/forgot-password', '/reset-password', '/api/auth/login']

export function middleware(request: NextRequest) {
  const token = request.cookies.get('auth_token')?.value
  const isPublic = publicRoutes.some((route) =>
    request.nextUrl.pathname.startsWith(route)
  )

  if (!token && !isPublic) {
    return NextResponse.redirect(new URL('/login', request.url))
  }

  return NextResponse.next()
}

export const config = {
  matcher: ['/((?!_next/static|_next/image|favicon.ico|.*\\.).*)'],
}
