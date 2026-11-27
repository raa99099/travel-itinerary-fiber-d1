import { useEffect, useState } from 'react';
import { api, getApiBaseUrl, getStoredToken, setStoredToken } from './services/api';
import AuthPanel from './components/AuthPanel';
import ItinerarySection from './components/ItinerarySection';
import SpotSection from './components/SpotSection';
import SpotNoteSection from './components/SpotNoteSection';

export default function App() {
  const [token, setToken] = useState(() => getStoredToken());
  const [healthStatus, setHealthStatus] = useState('Mengecek...');
  const [selectedItineraryId, setSelectedItineraryId] = useState('');
  const [selectedSpotId, setSelectedSpotId] = useState('');
  const [refreshKey, setRefreshKey] = useState(0);
  const [banner, setBanner] = useState(null);

  function handleNotify(type, message) {
    setBanner({ type, message });
  }

  function handleTokenChange(newToken) {
    setToken(newToken);
    setStoredToken(newToken);
    setRefreshKey((k) => k + 1);
  }

  async function checkHealth() {
    setHealthStatus('Mengecek...');
    try {
      await api.health();
      setHealthStatus('Online');
    } catch {
      setHealthStatus('Offline / Perlu CORS');
    }
  }

  useEffect(() => {
    checkHealth();
  }, []);

  return (
    <div className="app-shell">
      <header className="topbar">
        <div>
          <span className="eyebrow">Full-Stack Workspace</span>
          <h1>Travel Itinerary</h1>
          <p className="muted">
            API Base URL: <code>{getApiBaseUrl()}</code>
          </p>
        </div>
        <div className="topbar-right">
          <div className="step-pill">
            <span>Itinerary: <strong>{selectedItineraryId ? `#${selectedItineraryId}` : '-'}</strong></span>
            <span>•</span>
            <span>Spot: <strong>{selectedSpotId ? `#${selectedSpotId}` : '-'}</strong></span>
          </div>
          <div className="health-box">
            <span className={healthStatus === 'Online' ? 'dot dot-ok' : 'dot dot-warn'} />
            <span>Backend: <strong>{healthStatus}</strong></span>
            <button type="button" className="btn btn-ghost btn-sm" onClick={checkHealth}>
              Cek Ulang
            </button>
          </div>
        </div>
      </header>

      {banner && (
        <div className={`alert alert-${banner.type}`}>
          <span>{banner.message}</span>
          <button type="button" className="close-btn" onClick={() => setBanner(null)}>
            ×
          </button>
        </div>
      )}

      <main className="workspace-grid">
        <AuthPanel token={token} onTokenChange={handleTokenChange} onNotify={handleNotify} />
        <ItinerarySection
          selectedItineraryId={selectedItineraryId}
          onSelectItinerary={(id) => {
            setSelectedItineraryId(id);
            setSelectedSpotId('');
          }}
          onNotify={handleNotify}
          refreshKey={refreshKey}
        />
        <SpotSection
          itineraryId={selectedItineraryId}
          selectedSpotId={selectedSpotId}
          onSelectSpot={setSelectedSpotId}
          onNotify={handleNotify}
          refreshKey={refreshKey}
        />
        <SpotNoteSection spotId={selectedSpotId} onNotify={handleNotify} refreshKey={refreshKey} />
      </main>
    </div>
  );
}
