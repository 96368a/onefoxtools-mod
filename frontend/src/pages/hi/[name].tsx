export default function Hi() {
  const navigate = useNavigate()
  const params = useParams()

  return (
    <div>
      <div class="i-carbon-pedestrian text-4xl inline-block" />
      <p>
        Hi, {params.name}
      </p>
      <p class="text-sm op50">
        <em>Dynamic route!</em>
      </p>

      <div>
        <button
          class="btn m-3 text-sm mt-8"
          onClick={() => navigate(-1)}
        >
          Back
        </button>
      </div>
    </div>
  )
}
