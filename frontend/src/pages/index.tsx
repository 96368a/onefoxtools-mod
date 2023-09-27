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
      console.log(result)
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
    <div class=' pt-4'>
      <Search configs={configs} show={showSearch} setShow={setShowSearch} />
      <div class="card w-full mx-2 bg-base-100 shadow-xl">
        <div class="card-body">
          <div class="flex justify-center gap-2">
            <input type="search" class='input input-bordered w-full max-w-xs' onfocus={() => setShowSearch(true)} />
            <div class='flex gap-2'>
              <button class='btn btn-success'>搜索</button>
              <button class='btn btn-warning' onclick={refresh}>刷新</button>
            </div>
          </div>
        </div>
      </div>

      {/* <button class="btn" onclick={t}>233</button> */}
      <For each={configs}>{type => (
        <div class="card w-full my-2 bg-base-100 shadow-xl mx-2">
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
