import { useState } from "react";

function App() {
  const fetchData = () => {
    fetch(`http://localhost:${import.meta.env.VITE_PORT}/`)
      .then((response) => response.text())
      .then((data) => setMessage(data))
      .catch((error) => console.error("Error fetching data:", error));
  };
  const [message, setMessage] = useState<string>("");

  return (
    <>
      <h1 className="text-3xl font-bold underline">
        Welcome to the URL Crawler Application
      </h1>
      <button
        className="bg-blue-500 text-white p-2 rounded-md"
        onClick={fetchData}
      >
        Click to fetch from Go server
      </button>
      {message && (
        <div>
          <h2>Server Response:</h2>
          <p>{message}</p>
        </div>
      )}
    </>
  );
}

export default App;
