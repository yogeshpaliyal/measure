export default function ANRs({ params }: { params: { teamId: string } }) {
  return (
    <main className="flex flex-col items-center justify-between selection:bg-yellow-200/75">
      <div className="flex flex-col items-center md:w-4/5 2xl:w-3/5 px-16">
        <div className="py-24" />
        <p>ANRs</p>
        <div className="py-24" />
      </div>
    </main>
  )
}
