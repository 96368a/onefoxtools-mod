import { Router, useRoutes } from '@solidjs/router'
import { onMount } from 'solid-js'
import { themeChange } from 'theme-change'
import toast, { Toaster } from 'solid-toast'
import { KeepAliveProvider } from 'solid-keep-alive'
import { EventsOn } from 'wailsjs/runtime/runtime'
import Cookie from 'js-cookie'
import Footer from './components/Footer'
import Dialog from './components/Dialog'
import routes from '~solid-pages'

export default function App() {
  const Routes = useRoutes(routes)
  onMount(async () => {
    if (!Cookie.get('toastEvent')) {
      EventsOn('toast.success', (msg: string) => {
        toast.success(msg)
      })
      EventsOn('toast.error', (msg: string) => {
        toast.error(msg)
      })
      Cookie.set('toastEvent', 'true')
    }
    themeChange()
  })
  return (
    <main class="h-screen flex flex-col bg-base-200 text-center font-sans text-gray-700 dark:text-gray-200">
      <Toaster position='top-left' />
      <KeepAliveProvider>
        <Router>
          <Routes />
          <Footer />
        </Router>
      </KeepAliveProvider>
      <Dialog />
    </main>
  )
}
