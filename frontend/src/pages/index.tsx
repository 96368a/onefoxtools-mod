// import { Greet } from 'wailsjs/go/main/App'
import type { main } from 'wailsjs/go/models'

import toast from 'solid-toast'
import { Start } from '../../wailsjs/go/main/CONFIG'
import Search from '~/components/SearchUI'
import DataStore from '~/store/data'

export default function () {
  return (
    <Index></Index>
  )
}

function Index() {
  const { configs, refresConfig } = DataStore
  const [showSearch, setShowSearch] = createSignal(false)
  const navigate = useNavigate()

  const [time, setTime] = createSignal(-1)
  // 刷新数据
  async function refresh() {
    // 记录开始时间
    toast.promise(
      (async function () {
        const startTime = new Date().getTime()
        await refresConfig()
        const endTime = new Date().getTime()
        setTime((endTime - startTime) / 1000)
      }()),
      {
        loading: '加载中...',
        success: () => <span>加载完成,耗时{time()}秒</span>,
        error: '加载出错',
      },
    )
  }

  function start(c: main.Config) {
    toast.promise(
      (async function () {
        const startTime = new Date().getTime()
        await Start(c)
        const endTime = new Date().getTime()
        setTime((endTime - startTime) / 1000)
      }()),
      {
        loading: `启动 ${c.name} 中...`,
        success: () => <span> {c.name} 启动成功</span>,
        error: `${c.name} 执行出错`,
      },
    )
  }

  return (
    <div class=''>
      <Search configs={configs} show={showSearch} setShow={setShowSearch} />
      <div class="px-2 pt-2">
        <div class="w-full bg-base-100 shadow-xl card">
          <div class="card-body">
            <div class='absolute right-2 top-2 flex gap-1'>
              <button class='h-8 w-8 rounded-full btn-ghost' onclick={() => navigate('/setting/env')}>
                <div class='i-carbon-settings mx-auto w-6'></div>
              </button>
            </div>
            <div class='justify-center py-2 text-2xl card-title'>
              <h1>末影工具箱</h1>
            </div>
            <div class="flex justify-center gap-2">
              <input type="search" id="search" class='max-w-xs w-full input input-bordered input-sm' onfocus={() => setShowSearch(true)} />
              <div class='flex gap-2'>
                <button class='px-10 btn btn-success btn-sm' onclick={() => setShowSearch(true)}>搜索</button>
                <button class='px-10 btn btn-warning btn-sm' onclick={refresh}>刷新</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* <button class="btn" onclick={t}>233</button> */}
      <For each={configs}>{type => (
        <div class="px-3">
          <div class="my-2 w-full bg-base-100 shadow-xl card">
            <div class="items-center text-center card-body">
              <h2 class="card-title">{type.type}</h2>
              <div class='w-full'>
                <div class="card-actions">
                  <For each={type.config}>{
                    c => (
                      <button class="btn btn-outline btn-sm" onclick={() => start(c)}>{c.name}</button>
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
