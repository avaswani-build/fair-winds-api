import { useState } from 'react'
import './App.css'

function App() {
  const [lat, setLat] = useState("");
  const [long, setLong] = useState("");

  const handleSubmit = async () => {
    console.log(lat, long);
  }

  return (
    <>
      <section> 
        <div>
          <h1>Fair Winds</h1>

          <input
            placeholder="Latitude"
            value={lat}
            onChange={(e) => setLat(e.target.value)}
          />

          <input
          placeholder="Longitude"
          value={long}
          onChange={(e) => setLong(e.target.value)}
          />

          <button onClick={handleSubmit}>
            Check Conditions
          </button>
        </div>
      </section>
    </>
  );
}

export default App
