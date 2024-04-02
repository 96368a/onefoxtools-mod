import { GetBuildTime, GetGitCommit, GetVersion } from 'wailsjs/go/main/GOContext'

export default function () {
  const [version, setVersion] = createSignal('')
  const [buildTime, setBuildTime] = createSignal('')
  const [commit, setCommit] = createSignal('')
  onMount(async () => {
    GetVersion().then(setVersion)
    GetBuildTime().then(setBuildTime)
    GetGitCommit().then(setCommit)
  })

  return (
        <div class="h-full w-full bg-base-300 px-4 py-4 rounded-box">
            <div class="py-2">
                关于
            </div>
            <div>
                <div class="py-2">OnefoxTools-Mod</div>
                <div class="py-2">作者:木末君</div>
                <div class="py-2">Github: <a href="https://github.com/96368a/OnefoxTools-Mod">https://github.com/96368a/OnefoxTools-Mod</a></div>
                <div>当前版本: {version() || '未知版本'}</div>
                <div>编译时间: {buildTime() || '未知时间'}</div>
                <div>最后Commit版本: {commit() || '未知提交'}</div>
            </div>
        </div>
  )
}
