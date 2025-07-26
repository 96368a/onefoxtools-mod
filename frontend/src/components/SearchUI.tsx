import type { common } from 'wailsjs/go/models'
import PinyinMatch from 'pinyin-match'
import { Start } from 'wailsjs/go/main/GOContext'
import toast from 'solid-toast'

interface SearchProps {
  configs: common.TypeConfig[]
  show: () => boolean
  setShow: (b: boolean) => void
  showContextMenu: (e: MouseEvent, c: common.Config) => void
  hiddenContentMenu: () => void
}

export default function Search({ configs, show, setShow, showContextMenu, hiddenContentMenu }: SearchProps) {
  const [searchString, setSearchString] = createSignal('')
  const [searchResults, setSearchResults] = createSignal<common.Config[]>([])

  createEffect(() => {
    // 显示搜索ui时聚焦搜索框
    if (show()) {
      (document.querySelector('#root') as HTMLDivElement).onwheel = (e) => {
        e.preventDefault()
      }
      (document.querySelector('input[type=\'search\']') as HTMLInputElement).focus()
    }
    else {
      (document.querySelector('#root') as HTMLDivElement).onwheel = null
      hiddenContentMenu()
    }
  })

  function start(c: common.Config) {
    toast.promise(
      Start(c),
      {
        loading: `启动 ${c.name} 中...`,
        success: () => <span> {c.name} 启动成功</span>,
        error: `${c.name} 执行出错`,
      },
    )
  }

  function rs() {
    if (searchString() === '' || configs === undefined) {
      setSearchResults([] as common.Config[])
    }
    else {
      const results = []
      for (const type of configs) {
        for (const c of type.config) {
          if (PinyinMatch.match(c.name, searchString()))
            results.push(c)
        }
      }
      setSearchResults(results)
    }
  }

  return (
    <div>
      <Show when={show()}>
        <button class="absolute right-4 top-4 z-1001 btn btn-square btn-sm" onclick={() => setShow(false)}>
          ❌
        </button>
        <div class='fixed z-200 w-screen pt-20' onclick={() => setShow(false)}>
          <input type="search" class="max-w-lg w-full input input-bordered" placeholder="请输入搜索关键字，支持拼音缩写"
            maxlength="-1"
            value={searchString()} onKeyUp={e => setSearchString(e.currentTarget.value) && rs()} onclick={e => e.stopPropagation()} />

          <div class='p-4'>
            <div class='justify-center card-actions'>
              <For each={searchResults()}>
                {
                  c => (
                    <button class='btn btn-sm'
                    oncontextmenu={e => showContextMenu(e, c)}
                     onclick={(e) => {
                       start(c)
                       e.stopPropagation()
                     }}>
                      {c.name}
                    </button>
                  )
                }
              </For>
            </div>
          </div>
        </div>
        <div class="fixed z-100 h-screen w-screen bg-black opacity-50" onclick={() => setShow(false)}>
        </div>
      </Show>
    </div>
  )
}
