export default function LoadingSpinner() {
  return (
    <i
      className={`fas fa-spinner-third animate-spin-clean text-white text-3xl`}
    ></i>
  );
}

export function LoadingScreen({
  message = "Loading...",
}: {
  message?: string;
}) {
  return (
    <div className="min-h-screen bg-gray-900 flex flex-col items-center justify-center">
      <LoadingSpinner />
      <p className="mt-4 text-gray-400">{message}</p>
    </div>
  );
}
