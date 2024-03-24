import { Quit } from 'wailsjs/runtime/runtime'

export default function Init() {
  const params = useParams()
  const navigate = useNavigate()
  console.debug(params)
  function exitHandle() {
    Quit()
  }

  return (
        <div class="h-full w-full flex flex-1 justify-center py-20">
            <div class="w-xl bg-base-100 shadow-xl card">
                <div class="card-body">
                    <h2 class="justify-center card-title">{params.type && params.type === 'init' ? '提示' : '貌似出错了'}</h2>
                    <p>{params.msg ? params.msg : '未知错误'}</p>
                    <div class="justify-center gap-6 card-actions">
                        <Show
                            when={params.type && params.type === 'init'}
                            fallback={
                                <>
                                    <button class="btn btn-accent" onclick={() => navigate('/')}>重新加载</button>
                                    <button class="btn btn-error" onclick={exitHandle}>退出</button>
                                </>
                            }
                        >
                            <button class="btn btn-accent">初始化</button>
                            <button class="btn btn-error" onclick={exitHandle}>退出</button>
                        </Show>
                    </div>
                </div>
            </div>
        </div>
  )
}
