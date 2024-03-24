import toast from 'solid-toast'
import { GetStartTime } from 'wailsjs/go/main/GOContext'
import DataStore from '~/store/data'

export default () => {
  const navigate = useNavigate()
  // EventsOn('navigate', (path: string) => {
  //   navigate(path)
  // })
  const { refresConfig, refreshEnv } = DataStore
  onMount(async () => {
    refresConfig().then(async () => {
      await refreshEnv()
      GetStartTime().then((t) => {
        const startTime = new Date(t).getTime()
        const endTime = new Date().getTime()
        toast.success(`加载完成，耗时${(endTime - startTime) / 1000}秒`)
      })
      navigate('/home')
    }).catch((e: string) => {
      toast.error('加载配置文件出错')
      if (e.search('cannot find the file') !== -1) {
        navigate('/init', {
          state: {
            type: 'init',
            msg: '未检测到配置文件，点击开始初始化',
          },
        })
      }
      // location.href = '/init?type=init&msg=未检测到配置文件，点击开始初始化'
    })
  })
  return (
    <div class='flex flex-1 flex-col items-center justify-center gap-4'>
      <div class="icon-btn i-svg-spinners:tadpole text-6xl"></div>
      <div>加载中...</div>
    </div>)
}
