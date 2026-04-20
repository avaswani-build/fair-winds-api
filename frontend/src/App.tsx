import { getSummary } from "./api/summary";
import { useState } from 'react'
import type { SummaryResponse } from "./types/summary"
import './App.css'

function App() {
  const [lat, setLat] = useState("");
  const [long, setLong] = useState("");
  const [result, setResult] = useState<SummaryResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");

  const handleSubmit = async () => {
    setError("");
    setResult(null);

    const latNum = parseFloat(lat);
    const longNum = parseFloat(long);

    if (Number.isNaN(latNum) || Number.isNaN(longNum)) {
      setError("Please enter valid latitude and longitude values.");
      return;
    }

    try {
      setLoading(true);
      const summary = await getSummary(latNum, longNum);
      setResult(summary);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Something went wrong.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="app">
      <div className="card">
        <h1>Fair Winds</h1>
        <p className="subtitle">Check sailing conditions for a location.</p>

        <div className="input-group">
          <label htmlFor="lat">Latitude</label>
          <input
            id="lat"
            type="text"
            value={lat}
            onChange={(e) => setLat(e.target.value)}
            placeholder="e.g. 40.7128"
          />
        </div>

        <div className="input-group">
          <label htmlFor="lng">Longitude</label>
          <input
            id="lng"
            type="text"
            value={long}
            onChange={(e) => setLong(e.target.value)}
            placeholder="e.g. -74.0060"
          />
        </div>

        <button onClick={handleSubmit} disabled={loading}>
          {loading ? "Checking..." : "Check Conditions"}
        </button>

        {error && <p className="error">{error}</p>}

        {result && (
          <div className="result">
            <h2>Conditions Summary</h2>
            <p>
              <strong>Location:</strong> {result.location}
            </p>
            <p>
              <strong>Average Wind:</strong> {result.windAvgKts} kts
            </p>
            <p>
              <strong>Gusts:</strong> {result.gustKts} kts
            </p>
            <p>
              <strong>Wind Direction:</strong> {result.windDir}
            </p>
            <p>
              <strong>Status:</strong> {result.recommendation.status}
            </p>
            <p>
              <strong>Reason:</strong> {result.recommendation.reason}
            </p>
          </div>
        )}
      </div>
    </div>
  );
}

export default App
