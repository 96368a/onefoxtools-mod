import toast from 'solid-toast'
import { GenerateConfig } from 'wailsjs/go/main/GOContext'
import { Quit } from 'wailsjs/runtime/runtime'

export default function Init() {
  interface State {
    type: string
    msg: string
  }
  const location = useLocation<State>()
  const state = createMemo(() => location.state || {} as State)
  const navigate = useNavigate()
  const [isGenerating, setIsGenerating] = createSignal(false)
  function genConfig() {
    setIsGenerating(true)
    toast.promise(
      (GenerateConfig().then(() => {
        navigate('/')
        setIsGenerating(false)
      })),
      {
        loading: '正在生成配置文件...',
        success: '生成配置文件成功',
        error: '生成配置文件出错',
      },
    )
  }

  return (
    <div class="h-full w-full flex flex-1 justify-center py-20">
      <div class="w-xl bg-base-100 shadow-xl card">
        <div class="card-body">
          <h2 class="justify-center card-title">{state().type === 'init' ? '提示' : '貌似出错了'}</h2>
          <p>{state().msg || '未知错误'}</p>
          <div class="justify-center gap-6 card-actions">
            <Show
              when={state().type === 'init'}
              fallback={
                <>
                  <button class="btn btn-accent" onclick={() => navigate('/')}>重新加载</button>
                  <button class="btn btn-error" onclick={Quit}>退出</button>
                </>
              }
            >
              <Show when={!isGenerating()} fallback={
                <button class="btn btn-accent">
                  <span class="loading loading-spinner"></span>
                  生成中
                </button>
              }>
                <button class="btn btn-accent" onclick={genConfig} disabled={isGenerating()}>初始化</button>
              </Show>
              <button class="btn btn-error" onclick={Quit}>退出</button>
            </Show>
          </div>
        </div>
      </div>
    </div>
  )
}
