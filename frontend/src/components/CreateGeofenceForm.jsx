import { useState } from "react";

const API_URL = import.meta.env.VITE_API_URL;

function CreateGeofenceForm({ onClose, onCreated }) {
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [category, setCategory] = useState("delivery_zone");
  const [coordinates, setCoordinates] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

  const handleSubmit = async (e) => {
    e.preventDefault();

    setError("");
    setSuccess("");

    try {
      const parsedCoordinates = coordinates
        .trim()
        .split("\n")
        .map((line) => {
          const [lat, lng] = line.split(",").map(Number);

          if (Number.isNaN(lat) || Number.isNaN(lng)) {
            throw new Error(
              "Each coordinate must be in the format: latitude,longitude"
            );
          }

          return [lat, lng];
        });

      if (parsedCoordinates.length < 4) {
        throw new Error("At least 4 coordinates are required.");
      }

      setLoading(true);

      const response = await fetch(`${API_URL}/geofences`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name,
          description,
          category,
          coordinates: parsedCoordinates,
        }),
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.error || "Failed to create geofence");
      }

      setSuccess("Geofence created successfully.");

      if (onCreated) {
        onCreated(data);
      }

      setName("");
      setDescription("");
      setCoordinates("");
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="form-overlay">
      <div className="form-card">
        <div className="form-header">
          <h2>Create Geofence</h2>

          <button type="button" onClick={onClose}>
            ✕
          </button>
        </div>

        <form onSubmit={handleSubmit}>
          <label>
            Name
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="Downtown Delivery Zone"
              required
            />
          </label>

          <label>
            Description
            <textarea
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Main delivery area"
              rows={3}
            />
          </label>

          <label>
            Category
            <select
              value={category}
              onChange={(e) => setCategory(e.target.value)}
            >
              <option value="delivery_zone">Delivery Zone</option>
              <option value="restricted_zone">Restricted Zone</option>
              <option value="toll_zone">Toll Zone</option>
              <option value="customer_area">Customer Area</option>
            </select>
          </label>

          <label>
            Coordinates
            <textarea
              value={coordinates}
              onChange={(e) => setCoordinates(e.target.value)}
              placeholder={`37.7749,-122.4194
37.7849,-122.4194
37.7849,-122.4094
37.7749,-122.4094
37.7749,-122.4194`}
              rows={7}
              required
            />
          </label>

          <small>
            Enter one latitude,longitude pair per line. The first and last
            coordinates must be identical.
          </small>

          {error && <div className="form-error">{error}</div>}

          {success && <div className="form-success">{success}</div>}

          <div className="form-actions">
            <button type="button" onClick={onClose}>
              Cancel
            </button>

            <button type="submit" disabled={loading}>
              {loading ? "Creating..." : "Create Geofence"}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}

export default CreateGeofenceForm;