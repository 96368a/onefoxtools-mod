import type { common } from 'wailsjs/go/models'

import toast from 'solid-toast'
import Cookies from 'js-cookie'
import { GetStartTime, OpenFolderInExplorer, Start } from '../../wailsjs/go/main/GOContext'
import Search from '~/components/SearchUI'
import DataStore from '~/store/data'

export default function () {
  return (
    <Index></Index>
  )
}

function Index() {
  const { configs, envConfig, refresConfig } = DataStore
  const [showSearch, setShowSearch] = createSignal(false)
  const navigate = useNavigate()

  const [time, setTime] = createSignal(-1)
  // 刷新数据
  async function refresh() {
    // 记录开始时间
    toast.promise(
      (async function () {
        const startTime = new Date().getTime()
        await refresConfig().catch((e: string) => {
          console.error(e)
        })
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
  onMount(() => {
    if (Cookies.get('init') !== 'true') {
      GetStartTime().then((t) => {
        const startTime = new Date(t).getTime()
        const endTime = new Date().getTime()
        toast.success(`加载完成，耗时${(endTime - startTime) / 1000}秒`)
      })
      Cookies.set('init', 'true')
    }
  })

  function start(c: common.Config) {
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

  const [selectedEnv, setSelectedEnv] = createSignal({} as common.Config)
  let contentMenuRef = document.createElement('ul') // eslint-disable-line prefer-const
  const [isShowContentMenu, setShowContentMenu] = createSignal(false)
  createEffect(() => {
    if (isShowContentMenu())
      contentMenuRef.classList.add('z-999')
    else
      contentMenuRef.classList.remove('z-999')
  })
  // 计算元素到顶部的距离
  const CalcCoord = (element: EventTarget | null): number | null => {
    if (!(element instanceof HTMLElement))
      return null

    let actualTop = element.offsetTop
    let current = element.offsetParent as HTMLElement | null

    while (current !== null) {
      actualTop += current.offsetTop + current.clientTop
      current = current.offsetParent as HTMLElement | null
    }

    return actualTop + element.clientHeight - 2
  }

  const contentMentHandler = (e: MouseEvent, c: common.Config) => {
    e.preventDefault()
    e.stopPropagation()

    const target = e.currentTarget as HTMLElement
    const clientRect = target.getBoundingClientRect()
    // 设置菜单位置
    contentMenuRef.style.top = `${CalcCoord(target)}px`
    contentMenuRef.style.left = `${clientRect.x}px`

    setSelectedEnv(c)
    setShowContentMenu(true)
  }

  const hiddenContentMenu = () => {
    setShowContentMenu(false)
  }
  const openFolder = () => {
    // 使用正则表达式匹配绝对路径
    const regex = /cd (.+?) &&/i
    const match = selectedEnv().command.match(regex)
    if (match) {
      OpenFolderInExplorer(match[1]).then(() => {
        toast.success('打开文件夹成功')
        hiddenContentMenu()
      }).catch((e) => {
        toast.error(e)
      })
    }
    else {
      toast.error('未找到文件夹路径')
    }
  }

  return (
    <div onclick={hiddenContentMenu} oncontextmenu={hiddenContentMenu}>
      <Search configs={configs} show={showSearch} setShow={setShowSearch} />
      <div class="relative">
        <ul class="absolute right-0 top-12 z--1 w-56 bg-base-200 menu rounded-box"
          ref={contentMenuRef}>
          <li><a>查看信息</a></li>
          <li onclick={openFolder}><span>在文件夹中打开</span></li>
        </ul>
      </div>

      <div class="px-2 pt-2">
        <div class="w-full bg-base-100 shadow-xl card">
          <div class="card-body">
            <div class='absolute right-2 top-2 flex gap-1'>
              <button class='h-8 w-8 rounded-full btn-ghost' onclick={() => navigate('/setting/base')}>
                <div class='i-carbon-settings mx-auto w-6'></div>
              </button>
            </div>
            <div class='justify-center py-2 text-2xl card-title'>
              <h1>{envConfig.title ? envConfig.title : '安全工具箱'}</h1>
            </div>
            <div class="flex justify-center gap-2">
              <input type="search" id="search" class='max-w-xs w-full input input-bordered input-sm' onfocus={() => setShowSearch(true)} placeholder="请输入搜索关键字" />
              <div class='flex gap-2'>
                <button class='px-10 btn btn-success btn-sm' onclick={() => setShowSearch(true)}>搜索</button>
                <button class='px-10 btn btn-warning btn-sm' onclick={refresh}>刷新</button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <For each={configs}>{type => (
        <div class="px-3">
          <div class="my-2 w-full bg-base-100 shadow-xl card">
            <div class="items-center text-center card-body">
              <h2 class="card-title">{type.type}</h2>
              <div class='w-full'>
                <div class="card-actions">
                  <For each={type.config}>{
                    c => (
                      <button class="btn btn-outline btn-sm"
                        oncontextmenu={e => contentMentHandler(e, c)}
                        onclick={() => start(c)}>{c.name}</button>
                    )
                  }</For>
                </div>
              </div>
            </div>
          </div>
        </div>
      )
      }
      </For>

    </div>
  )
}
