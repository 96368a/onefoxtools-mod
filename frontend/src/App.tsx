import { Router, useRoutes } from '@solidjs/router'
import { onMount } from 'solid-js'
import { themeChange } from 'theme-change'
import Footer from './components/Footer'
import routes from '~solid-pages'

export default function App() {
  const Routes = useRoutes(routes)
  onMount(async () => {
    themeChange()
  })
  return (
    <main class="font-sans text-center text-gray-700 dark:text-gray-200">
      <Router>
        <Routes />
      </Router>
      <Footer />
    </main>
  )
}
