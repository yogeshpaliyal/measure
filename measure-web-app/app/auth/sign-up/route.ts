import { createClient } from '@/utils/supabase/server'
import { AuthApiError } from '@supabase/supabase-js'
import { NextResponse } from 'next/server'

export const dynamic = 'force-dynamic'

export async function POST(request: Request) {
  const requestUrl = new URL(request.url)
  const formData = await request.formData()
  const email = String(formData.get('email'))
  const supabase = createClient()

  const { error } = await supabase.auth.signInWithOtp({
    email,
    options: {
      emailRedirectTo: `${requestUrl.origin}/auth/callback`,
    },
  })

  if (error && error instanceof AuthApiError && error.status === 429) {
    console.log(error)
    return NextResponse.redirect(
      `${requestUrl.origin}/auth/login?error=Too many attempts, please try again in a minute.`,
      {
        // a 301 status is required to redirect from a POST to a GET route
        status: 301,
      }
    )
  }

  if (error) {
    console.log({ error })
    return NextResponse.redirect(
      `${requestUrl.origin}/auth/login?error=Could not authenticate`,
      {
        // a 301 status is required to redirect from a POST to a GET route
        status: 301,
      }
    )
  }

  return NextResponse.redirect(
    `${requestUrl.origin}/auth/login?message=Check email to continue to sign in`,
    {
      // a 301 status is required to redirect from a POST to a GET route
      status: 301,
    }
  )
}