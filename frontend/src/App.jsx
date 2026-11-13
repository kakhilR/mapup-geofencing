import { useEffect, useState } from "react";
import { MapContainer, TileLayer, Marker, Popup, Polygon } from "react-leaflet";
import L from "leaflet";
import "./App.css";
const API_URL = import.meta.env.VITE_API_URL;
const WS_URL = import.meta.env.VITE_WS_URL;


import CreateGeofenceForm from "./components/CreateGeofenceForm";

import markerIcon from "leaflet/dist/images/marker-icon.png";
import markerShadow from "leaflet/dist/images/marker-shadow.png";

const defaultIcon = L.icon({
  iconUrl: markerIcon,
  shadowUrl: markerShadow,
  iconSize: [25, 41],
  iconAnchor: [12, 41],
});

L.Marker.prototype.options.icon = defaultIcon;

async function fetchJSON(url) {
  const res = await fetch(url);
  const contentType = res.headers.get("content-type") || "";

  if (!res.ok) {
    const errorBody = contentType.includes("application/json")
      ? JSON.stringify(await res.json())
      : await res.text();
    throw new Error(errorBody || `Request failed with status ${res.status}`);
  }

  if (!contentType.includes("application/json")) {
    throw new Error("Expected JSON response");
  }

  return res.json();
}

function App() {
  const [connected, setConnected] = useState(false);
  const [alerts, setAlerts] = useState([]);
  const [geofences, setGeofences] = useState([]);
  const [vehicleLocation, setVehicleLocation] = useState(null);
  const [showGeofenceForm, setShowGeofenceForm] = useState(false);

  useEffect(() => {
    const fetchGeofences = async () => {
      try {
        const response = await fetch(`${API_URL}/geofences`);

        console.log("Geofence response status:", response.status);
        console.log(
          "Geofence response headers:",
          [...response.headers.entries()]
        );

        const text = await response.text();

        console.log("Geofence raw response:", text);

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${text}`);
        }

        const data = JSON.parse(text);

        console.log("Geofence parsed data:", data);

        setGeofences(data.geofences || []);
      } catch (error) {
        console.error("Failed to fetch geofences:", error);
      }
    };

    fetchGeofences();
  }, []);
  
  useEffect(() => {
    fetchJSON(
      `${API_URL}/vehicles/location/d51d51d7-caef-4c94-898d-b6b282cf07d9`
    )
      .then((data) => {
        console.log("Vehicle location:", data);
        setVehicleLocation(data);
      })
      .catch((err) => {
        console.error("Failed to fetch vehicle location:", err);
        setVehicleLocation(null);
      });
  }, []);

  useEffect(() => {
    const fetchVehicleLocation = async () => {
      try {
        const response = await fetch(
          `${API_URL}/vehicles/location/d51d51d7-caef-4c94-898d-b6b282cf07d9`
        );

        if (!response.ok) {
          throw new Error(`HTTP ${response.status}`);
        }

        const data = await response.json();
        setVehicleLocation(data);
      } catch (error) {
        console.error("Failed to fetch vehicle location:", error);
      }
    };

    fetchVehicleLocation();

    const interval = setInterval(fetchVehicleLocation, 3000);

    return () => clearInterval(interval);
  }, []);

  useEffect(() => {
    const ws = new WebSocket(`${WS_URL}/ws/alerts`);

    ws.onopen = () => {
      console.log("WebSocket connected");
      setConnected(true);
    };

    ws.onmessage = (event) => {
      const alert = JSON.parse(event.data);

      setAlerts((prev) => [alert, ...prev]);
    };

    ws.onclose = () => {
      console.log("WebSocket disconnected");
      setConnected(false);
    };

    ws.onerror = (error) => {
      console.error("WebSocket error:", error);
      setConnected(false);
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="app">
      <header className="header">
        <div>
          <h1>MapUp Geofence Monitor</h1>
          <p>Real-time vehicle monitoring</p>
        </div>

        <div className={`connection ${connected ? "online" : "offline"}`}>
          <span className="dot"></span>
          {connected ? "Connected" : "Disconnected"}
        </div>
      </header>

      <main className="dashboard">
       <section className="map-container">
          <MapContainer
            center={[37.7799, -122.4144]}
            zoom={14}
            style={{ height: "100%", width: "100%" }}
          >
            <TileLayer
              attribution='&copy; OpenStreetMap contributors'
              url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png"
            />

            {vehicleLocation && (
              <Marker
                position={[
                  vehicleLocation.latitude,
                  vehicleLocation.longitude,
                ]}
              >
                <Popup>
                  <strong>Vehicle</strong>
                  <br />
                  ID: {vehicleLocation.vehicle_id?.slice(0, 8) || "unknown"}
                  <br />
                  Latitude: {vehicleLocation.latitude}
                  <br />
                  Longitude: {vehicleLocation.longitude}
                </Popup>
              </Marker>
            )}

              {geofences.map((geofence) => (
                <Polygon
                  key={geofence.id}
                  positions={geofence.coordinates}
                >
                  <Popup>
                    <strong>{geofence.name}</strong>
                    <br />
                    Category: {geofence.category}
                    <br />
                    {geofence.description}
                  </Popup>
                </Polygon>
              ))}
          </MapContainer>
        </section>

        <section className="alerts">
          <div className="alerts-header">
            <h2>Live Alerts</h2>
            <span>{alerts.length}</span>
          </div>

          {alerts.length === 0 ? (
            <div className="empty">
              Waiting for vehicle events...
            </div>
          ) : (
            <div className="alert-list">
              {alerts.map((alert) => (
                <div
                  key={alert.id}
                  className={`alert ${alert.event_type.toLowerCase()}`}
                >
                  <div className="alert-type">
                    {alert.event_type}
                  </div>

                  <div className="alert-message">
                    {alert.message}
                  </div>

                  <div className="alert-time">
                    {new Date(alert.created_at).toLocaleTimeString()}
                  </div>
                </div>
              ))}
            </div>
          )}

          <button onClick={() => setShowGeofenceForm(true)}>
            Create Geofence
          </button>

          {showGeofenceForm && (
            <CreateGeofenceForm
              onClose={() => setShowGeofenceForm(false)}
              onCreated={async () => {
                try {
                  const data = await fetchJSON(`${API_URL}/geofences`);
                  setGeofences(data.geofences || []);
                } catch (error) {
                  console.error("Failed to refresh geofences:", error);
                }
              }}
            />
          )}
        </section>
      </main>
    </div>
  );
}

export default App;