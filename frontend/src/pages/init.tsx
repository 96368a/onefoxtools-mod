import { Quit } from 'wailsjs/runtime/runtime'

export default function Init() {
  const [searchParams, setSearchParams] = useSearchParams()
  function exitHandle() {
    Quit()
  }

  return (
    <div class="h-full w-full flex flex-1 justify-center py-20">
      <div class="w-xl bg-base-100 shadow-xl card">
        <div class="card-body">
          <h2 class="justify-center card-title">{searchParams.type && searchParams.type === 'init' ? '提示' : '貌似出错了'}</h2>
          <p>{searchParams.msg ? searchParams.msg : '未知错误'}</p>
          <div class="justify-center gap-6 card-actions">
            <Show
              when={searchParams.type && searchParams.type === 'init'}
              fallback={
                <>
              <button class="btn btn-accent">重新加载</button>
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
