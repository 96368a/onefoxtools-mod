import { Quit } from 'wailsjs/runtime/runtime'
import type { ErrorMsg } from '~/types/ErrorMsg'

export default () => {
  const location = useLocation<ErrorMsg>()
  const navigate = useNavigate()
  function toIndex() {
    navigate('/')
  }

  return <div class="h-full w-full py-20">
        <div class="mx-auto w-96 bg-base-100 shadow-xl card">
            <div class="items-center text-center card-body">
                <h2 class="card-title">出现错误!</h2>
                <p class="py-4">{location.state?.msg}</p>
                <div class="justify-center card-actions">
                    <button class="btn btn-info" onclick={toIndex}>返回重试</button>
                    <button class="btn btn-error" onclick={Quit}>退出程序</button>
                </div>
            </div>
        </div>
    </div>
}
