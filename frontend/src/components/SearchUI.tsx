import type { common } from 'wailsjs/go/models'
import PinyinMatch from 'pinyin-match'
import { Start } from 'wailsjs/go/main/GOContext'
import toast from 'solid-toast'

export default function Search({ configs, show, setShow }: { configs: common.TypeConfig[]; show: () => boolean; setShow: (b: boolean) => void }) {
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
        <div class='fixed z-1000 w-screen pt-20' onclick={() => setShow(false)}>
          <input type="search" class="max-w-lg w-full input input-bordered" placeholder="请输入搜索关键字，支持拼音缩写"
            maxlength="-1"
            value={searchString()} onKeyUp={e => setSearchString(e.currentTarget.value) && rs()} onclick={e => e.stopPropagation()} />
          {/* <input type="search" class="w-200 px-4 py-2 rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:border-transparent" /> */}

          <div py-4>
            <ul class='flex flex-col gap-2'>
              <For each={searchResults()}>
                {
                  c => (
                    <li class='mx-auto w-80 cursor-pointer rounded py-1' onclick={(e) => {
                      start(c)
                      e.stopPropagation()
                    }}>
                      <button class="border btn btn-sm">{c.name}</button>
                    </li>
                  )
                }
              </For>
            </ul>
          </div>
        </div>
        <div class="fixed z-100 h-screen w-screen bg-black opacity-50" onclick={() => setShow(false)}>
        </div>
      </Show>
    </div>
  )
}
