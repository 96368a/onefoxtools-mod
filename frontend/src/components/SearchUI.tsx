import type { main } from 'wailsjs/go/models'
import PinyinMatch from 'pinyin-match'
import { Start } from 'wailsjs/go/main/CONFIG'

export default function Search({ configs, show, setShow }: { configs: main.TypeConfig[]; show: () => boolean; setShow: (b: boolean) => void }) {
  const [searchString, setSearchString] = createSignal('')
  const [searchResults, setSearchResults] = createSignal<main.Config[]>([])

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

  function rs() {
    if (searchString() === '' || configs === undefined) {
      setSearchResults([] as main.Config[])
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
          <div class='fixed z-1000 w-screen pt-20' onclick={() => setShow(false)}>
          <input type="search" class="max-w-lg w-full input input-bordered"placeholder="Search..."
          maxlength="-1"
          value={searchString()} onKeyUp={e => setSearchString(e.currentTarget.value) && rs()} onclick={e => e.stopPropagation()}/>
            {/* <input type="search" class="w-200 px-4 py-2 rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:border-transparent" /> */}

            <div py-4>
              <ul class='flex flex-col gap-2'>
                <For each={searchResults()}>
                  {
                    c => (
                      <li class='mx-auto w-80 cursor-pointer rounded py-1' onclick={(e) => {
                        Start(c)
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
