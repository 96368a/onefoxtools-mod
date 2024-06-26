import toast from 'solid-toast'

export default function () {
  const navigate = useNavigate()
  function not() {
    toast.error('暂未完善')
  }
  return (
        <div class="w-full flex gap-2 p-2">
            <div class="h-full bg-base-300 pt-4 rounded-box">
                <div>
                    设置
                </div>
                <ul class="menu">
                    <li onclick={() => navigate('/setting/base')} ondblclick={() => navigate('/setting/debug')}>
                        <a class="w-50 flex justify-start">
                            <div class='i-carbon-settings h-5 w-5'></div>
                            <span>基础设置</span>
                        </a>
                    </li>
                    <li onclick={() => navigate('/setting/env')} ondblclick={() => navigate('/setting/test')}>
                        <a class="flex justify-start">
                            <div class='i-carbon-tool-box h-5 w-5'></div>
                            <div>环境变量配置</div>
                        </a>
                    </li>
                    <li onclick={not} ondblclick={() => navigate('/setting/tools')}>
                        <a class="flex justify-start">
                            <div class='i-carbon-tool-kit h-5 w-5'></div>
                            <span>工具设置</span>
                        </a>
                    </li>
                    <li onclick={() => navigate('/setting/about')}>
                        <a class="flex justify-start">
                            <div class='i-carbon-information h-5 w-5'></div>
                            <span>关于</span>
                        </a>
                    </li>
                    <li onclick={() => navigate('/')}>
                        <a class="flex justify-start">
                            <div class='i-carbon-arrow-left h-5 w-5'></div>
                            <span>返回主页</span>
                        </a>
                    </li>
                </ul>
            </div>
            <Outlet />
        </div>
  )
}
