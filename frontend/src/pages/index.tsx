import toast from 'solid-toast'
import DataStore from '~/store/data'

export default () => {
  const navigate = useNavigate()
  // EventsOn('navigate', (path: string) => {
  //   navigate(path)
  // })
  const { refresConfig } = DataStore
  onMount(async () => {
    refresConfig().then(() => {
      navigate('/home')
    }).catch((e: Error) => {
      toast.error('加载配置文件出错')
      if (e.message.search('cannot find the file') !== -1) {
        navigate('/init', {
          state: {
            type: 'init',
            msg: '未检测到配置文件，点击开始初始化',
          },
        })
      }
      console.error(e)
    })
  })
  return (
    <div class='flex flex-1 flex-col items-center justify-center gap-4'>
      <div class="icon-btn i-svg-spinners:tadpole text-6xl"></div>
      <div>加载中...</div>
    </div>)
}
