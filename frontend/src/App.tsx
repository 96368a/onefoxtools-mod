import { Router, useRoutes } from '@solidjs/router'
import { onMount } from 'solid-js'
import { themeChange } from 'theme-change'
import toast, { Toaster } from 'solid-toast'
import { KeepAliveProvider } from 'solid-keep-alive'
import { GetStartTime } from 'wailsjs/go/main/GOContext'
import { EventsOn } from 'wailsjs/runtime/runtime'
import Footer from './components/Footer'
import Dialog from './components/Dialog'
import routes from '~solid-pages'
import DataStore from '~/store/data'

export default function App() {
  const Routes = useRoutes(routes)
  const { refresConfig, refreshEnv } = DataStore
  onMount(async () => {
    EventsOn('toast.success', (msg: string) => {
      toast.success(msg)
    })
    EventsOn('toast.error', (msg: string) => {
      toast.error(msg)
    })
    EventsOn('navigate', (path: string) => {
      window.location.href = path
    })
    themeChange()
    refresConfig().then(async () => {
      await refreshEnv()
      GetStartTime().then((t) => {
        const startTime = new Date(t).getTime()
        const endTime = new Date().getTime()
        toast.success(`加载完成，耗时${(endTime - startTime) / 1000}秒`)
      })
    }).catch(() => {
      toast.error('加载配置文件出错')
    })
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
