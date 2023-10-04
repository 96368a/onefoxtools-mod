import { Router, useRoutes } from '@solidjs/router'
import { onMount } from 'solid-js'
import { themeChange } from 'theme-change'
import { Toaster } from 'solid-toast'
import Footer from './components/Footer'
import routes from '~solid-pages'

export default function App() {
  const Routes = useRoutes(routes)

  onMount(async () => {
    themeChange()
  })
  return (
    <main class="bg-base-200 text-center font-sans text-gray-700 dark:text-gray-200">
      <Toaster position='top-left'/>
      <Router>
        <Routes />
      </Router>
      <Footer />
    </main>
  )
}
