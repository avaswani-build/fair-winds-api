import { getSummary } from "./api/summary";
import { useState } from 'react'
import { fetchTimeline } from "./api/timeline";
import type { TimelinePoint } from "./types/timeline";
import type { SummaryResponse } from "./types/summary"
import './App.css'

function formatHour(time: string) {
  return new Date(time).toLocaleTimeString([], {
    hour: "numeric",
    hour12: true,
  });
}

function formatDay(time: string) {
  return new Date(time).toLocaleDateString([], {
    weekday: "short",
    month: "short",
    day: "numeric",
  });
}

function levelLabel(level: string) {
  switch (level) {
    case "light":
      return "Light";
    case "medium":
      return "Medium";
    case "heavy":
      return "Heavy";
    default:
      return level;
  }
}

function App() {
  const [lat, setLat] = useState("");
  const [long, setLong] = useState("");
  const [result, setResult] = useState<SummaryResponse | null>(null);
  const [timeline, setTimeline] = useState<TimelinePoint[]>([]);
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
      const [summary, timelineData] = await Promise.all([
        getSummary(latNum, longNum),
        fetchTimeline(latNum, longNum),
      ]);

      setResult(summary);
      setTimeline(timelineData);
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

        <button
          onClick={() => {
            if (!navigator.geolocation) {
              setError("Geolocation not supported");
              return;
            }

            setLoading(true);
            navigator.geolocation.getCurrentPosition(
              async (pos) => {
                const latVal = pos.coords.latitude;
                const lngVal = pos.coords.longitude;

                setLat(latVal.toString());
                setLong(lngVal.toString());

                try {
                  const [summary, timelineData] = await Promise.all([
                    getSummary(latVal, lngVal),
                    fetchTimeline(latVal, lngVal),
                  ]);

                  setResult(summary);
                  setTimeline(timelineData);
                } catch (err) {
                  console.error("Fetch error:", err);
                  setError("Failed to fetch summary");
                } finally {
                  setLoading(false);
                }
              },
              (err) => {
                console.error("Geolocation error:", err);
                setError(`Geolocation error: ${err.message}`);
                setLoading(false);
              }
            );
          }}
        >
          Use My Location
        </button>

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

        {timeline.length > 0 && (
          <section className="timeline-card">
            <div className="timeline-header">
              <div>
                <p className="eyebrow">Wind outlook</p>
                <h2>Next sailing window</h2>
              </div>
              <span className="timeline-count">{timeline.length} hrs</span>
            </div>

            <div className="timeline-strip">
              {timeline.slice(0, 24).map((point) => (
                <div className="timeline-point" key={point.time}>
                  <div className={`wind-dot ${point.level}`} />
                  <span className="timeline-hour">{formatHour(point.time)}</span>
                </div>
              ))}
            </div>

            <div className="timeline-summary">
              <span>Now: {levelLabel(timeline[0].level)} wind</span>
              <span>{formatDay(timeline[0].time)}</span>
            </div>
          </section>
        )}
      </div>
    </div>
  );
}

export default App
