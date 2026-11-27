import { useEffect, useState } from 'react';
import { api, extractList } from '../services/api';

const EMPTY_FORM = {
  name: '',
  location: '',
  day_number: 1,
  notes: '',
};

export default function SpotSection({
  itineraryId,
  selectedSpotId,
  onSelectSpot,
  onNotify,
  refreshKey,
}) {
  const [spots, setSpots] = useState([]);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(false);

  async function fetchSpots() {
    setLoading(true);
    try {
      const res = await api.listSpots(itineraryId);
      const all = extractList(res);
      const filtered = itineraryId
        ? all.filter(
            (s) =>
              s.itinerary_id === undefined ||
              s.itinerary_id === null ||
              String(s.itinerary_id) === String(itineraryId)
          )
        : all;
      setSpots(filtered);
      if (filtered.length > 0 && !filtered.some((s) => String(s.id) === String(selectedSpotId))) {
        onSelectSpot(filtered[0].id);
      } else if (filtered.length === 0) {
        onSelectSpot('');
      }
    } catch (err) {
      onNotify('error', `Spots: ${err.message}`);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchSpots();
  }, [itineraryId, refreshKey]);

  function startEdit(spot) {
    setEditingId(spot.id);
    setForm({
      name: spot.name || '',
      location: spot.location || '',
      day_number: spot.day_number || 1,
      notes: spot.notes || '',
    });
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    if (!itineraryId) {
      onNotify('error', 'Pilih atau buat Itinerary terlebih dahulu.');
      return;
    }
    setLoading(true);
    try {
      const payload = {
        itinerary_id: Number(itineraryId),
        name: form.name.trim(),
        location: form.location.trim(),
        day_number: Number(form.day_number) || 1,
        notes: form.notes.trim(),
      };
      if (editingId) {
        await api.updateSpot(editingId, payload);
        onNotify('success', 'Spot berhasil diperbarui.');
      } else {
        const res = await api.createSpot(payload);
        const createdId = res?.data?.id || res?.id;
        if (createdId) onSelectSpot(createdId);
        onNotify('success', 'Spot baru ditambahkan.');
      }
      cancelEdit();
      await fetchSpots();
    } catch (err) {
      onNotify('error', err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id) {
    setLoading(true);
    try {
      await api.deleteSpot(id);
      if (String(selectedSpotId) === String(id)) onSelectSpot('');
      onNotify('success', 'Spot dihapus.');
      await fetchSpots();
    } catch (err) {
      onNotify('error', err.message);
    } finally {
      setLoading(false);
    }
  }

  return (
    <section className="card">
      <div className="card-header">
        <div>
          <h2>3. Spots (Destinasi)</h2>
          <p className="muted">
            Itinerary aktif: <strong>{itineraryId ? `#${itineraryId}` : 'Belum dipilih'}</strong>
          </p>
        </div>
        <button type="button" className="btn btn-ghost btn-sm" onClick={fetchSpots}>
          Refresh
        </button>
      </div>

      <form onSubmit={handleSubmit} className="form-grid">
        <label>
          Nama Destinasi *
          <input
            type="text"
            placeholder="Contoh: Pantai Kuta"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            required
          />
        </label>
        <label>
          Lokasi
          <input
            type="text"
            placeholder="Contoh: Badung, Bali"
            value={form.location}
            onChange={(e) => setForm({ ...form, location: e.target.value })}
          />
        </label>
        <label>
          Hari Ke-
          <input
            type="number"
            min="1"
            value={form.day_number}
            onChange={(e) => setForm({ ...form, day_number: e.target.value })}
          />
        </label>
        <label>
          Catatan Singkat
          <input
            type="text"
            placeholder="Opsional"
            value={form.notes}
            onChange={(e) => setForm({ ...form, notes: e.target.value })}
          />
        </label>
        <div className="form-actions span-2">
          <button type="submit" className="btn btn-primary" disabled={loading || !itineraryId}>
            {editingId ? 'Simpan Spot' : 'Tambah Spot'}
          </button>
          {editingId && (
            <button type="button" className="btn btn-ghost" onClick={cancelEdit}>
              Batal
            </button>
          )}
        </div>
      </form>

      <div className="item-list">
        {spots.length === 0 ? (
          <p className="empty-state">Belum ada spot untuk itinerary ini.</p>
        ) : (
          spots.map((spot) => {
            const isSelected = String(spot.id) === String(selectedSpotId);
            return (
              <div key={spot.id} className={isSelected ? 'item-row item-selected' : 'item-row'}>
                <div onClick={() => onSelectSpot(spot.id)} className="clickable">
                  <strong>#{spot.id} — {spot.name}</strong>
                  <div className="muted">
                    {spot.location || 'Lokasi belum diisi'}
                    {spot.day_number ? ` • Hari ke-${spot.day_number}` : ''}
                  </div>
                </div>
                <div className="row-actions">
                  <button
                    type="button"
                    className="btn btn-secondary btn-sm"
                    onClick={() => onSelectSpot(spot.id)}
                  >
                    {isSelected ? 'Dipilih' : 'Catatan'}
                  </button>
                  <button
                    type="button"
                    className="btn btn-ghost btn-sm"
                    onClick={() => startEdit(spot)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className="btn btn-danger btn-sm"
                    onClick={() => handleDelete(spot.id)}
                  >
                    Hapus
                  </button>
                </div>
              </div>
            );
          })
        )}
      </div>
    </section>
  );
}
