// import { Greet } from 'wailsjs/go/main/App'
import type { main } from 'wailsjs/go/models'

import { GetConfigs, InitConfig, Start } from '../../wailsjs/go/main/CONFIG'
import Search from '~/components/SearchUI'

export default function Index() {
  const [configs, setConfigs] = createStore<main.TypeConfig[]>([])
  const [showSearch, setShowSearch] = createSignal(false)
  onMount(async () => {
    await setData()
  })
  async function setData() {
    GetConfigs().then((result) => {
      // 工具类别按照index进行排序
      result.sort((a, b) => {
        if (a.index === 0)
          a.index = 1e10
        if (b.index === 0)
          b.index = 1e10
        return a.index - b.index
      })
      // 类比里的工具按照index进行排序
      result = result.map((c) => {
        c.config.sort((a, b) => {
          if (a.index === 0)
            a.index = 1e10
          if (b.index === 0)
            b.index = 1e10
          return a.index - b.index
        })
        return c
      })
      setConfigs(result)
    })
  }
  async function refresh() {
    await InitConfig()
    await setData()
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
    <div class=''>
      <Search configs={configs} show={showSearch} setShow={setShowSearch} />
      <div class="px-2 pt-2">
      <div class="card w-full bg-base-100 shadow-xl">
        <div class="card-body">
          <div class='justify-center card-title text-2xl py-2'>
          <h1>ONE-FOX工具箱</h1>
          </div>
          <div class="flex justify-center gap-2">
            <input type="search" id="search" class='input input-bordered w-full max-w-xs input-sm' onfocus={() => setShowSearch(true)} />
            <div class='flex gap-2'>
              <button class='btn btn-sm btn-success px-10' onclick={() => setShowSearch(true)}>搜索</button>
              <button class='btn btn-sm btn-warning px-10' onclick={refresh}>刷新</button>
            </div>
          </div>
        </div>
      </div>
      </div>

      {/* <button class="btn" onclick={t}>233</button> */}
      <For each={configs}>{type => (
        <div class="px-3">
        <div class="card w-full my-2 bg-base-100 shadow-xl">
          <div class="card-body   items-center text-center">
            <h2 class="card-title">{type.type}</h2>
            <div class='w-full'>
              <div class="card-actions">
                <For each={type.config}>{
                  c => (
                    <button class="btn btn-sm btn-outline" onclick={() => start(c)}>{c.name}</button>
                  )
                }</For>
                {/* <button class="btn btn-primary">Buy Now</button> */}
              </div>
            </div>
          </div>
        </div>
        </div>
        // <div class='shadow rounded mb-1 py-2 mx-4'>
        //   <h2 class='w-full text-center'>{type.type}</h2>
        //   <div class='flex gap-2 justify-start p-2 flex-wrap'>
        //     <For each={type.config}>{
        //       c => (
        //         <button class="btn text-xs truncate" onclick={() => start(c)}>{c.name}</button>
        //       )
        //     }</For>
        //   </div>
        // </div>
      )
      }
      </For>

    </div>
  )
}
