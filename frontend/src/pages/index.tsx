// import { Greet } from 'wailsjs/go/main/App'
import type { main } from 'wailsjs/go/models'
import PinyinMatch from 'pinyin-match'
import { GetConfigs, Start } from '../../wailsjs/go/main/CONFIG'

function Search({ configs, show, setShow }: { configs: { [key: string]: main.Config[] }; show: () => boolean; setShow: (b: boolean) => void }) {
  const [searchString, setSearchString] = createSignal('')
  const [searchResults, setSearchResults] = createSignal<main.Config[]>([])

  createEffect(() => {
    // 显示搜索ui时聚焦搜索框
    if (show())
      (document.querySelector('input[type=\'search\']') as HTMLInputElement).focus()
  })

  function rs() {
    if (searchString() === '' || configs === undefined) {
      setSearchResults([] as main.Config[])
    }
    else {
      const results = []
      for (const type in configs) {
        for (const c of configs[type]) {
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
      <div class='w-screen fixed z-1000 pt-20' onclick={() => setShow(false)}>
          <input type="search" class="w-200 px-4 py-2 rounded-md border border-gray-300 focus:outline-none focus:ring-2 focus:border-transparent" placeholder="Search..."
          value={searchString()} oninput={e => setSearchString(e.currentTarget.value) && rs()} onclick={e => e.stopPropagation()}/>

          <div py-4>
            <ul class='flex flex-col gap-2'>
              <For each={searchResults()}>
                {
                  c => (
                    <li class='bg-light-200 w-80 mx-auto rounded py-1 cursor-pointer' onclick={(e) => {
                      Start(c)
                      e.stopPropagation()
                    }}>{c.name}</li>
                  )
                }
              </For>
            </ul>
          </div>
        </div>
      <div class="w-screen h-screen z-100 opacity-50 fixed bg-black" onclick={() => setShow(false)}>
      </div>
      </Show>
    </div>
  )
}

export default function Index() {
  const [configs, setConfigs] = createStore<main.TypeConfig[]>([])
  const [showSearch, setShowSearch] = createSignal(false)
  onMount(async () => {
    // await setData()
  })
  async function setData() {
    GetConfigs().then((result) => {
      result.sort((a, b) => {
        if (a.index === 0)
          a.index = 1e10
        if (b.index === 0)
          b.index = 1e10
        return a.index - b.index
      })
      setConfigs(result)
      console.log(result)
    })
  }
  function start(c: main.Config) {
    Start(c)
  }
  // function t() {
  //   GetConfigs().then((result) => {
  //     console.log(result)
  //     setConfigs(result)
  //   })
  // }

  return (
    <div>
      <Search configs={configs} show={showSearch} setShow={setShowSearch}/>
      <div class="py-2">
        <input type="search" onfocus={() => setShowSearch(true)} class='border rounded outline-none px-2 py-1'/>
        <button class='bg-gray-400 mx-2 px-6 text-light py-1 rounded hover:bg-gray-5'>搜索</button>
        <button class='bg-gray-400   px-6 text-light py-1 rounded hover:bg-gray-5' onclick={setData}>刷新</button>
      </div>
      {/* <button class="btn" onclick={t}>233</button> */}
      <For each={configs}>{type => (
        <div class='shadow rounded mb-1 py-2 mx-4'>
          <h2 class='w-full text-center'>{type.type}</h2>
          <div class='flex gap-2 justify-start p-2 flex-wrap'>
            <For each={type.config}>{
              c => (
                <button class="btn text-xs truncate" onclick={() => start(c)}>{c.name}</button>
              )
            }</For>
          </div>
        </div>
      )
      }
      </For>

    </div>
  )
}
