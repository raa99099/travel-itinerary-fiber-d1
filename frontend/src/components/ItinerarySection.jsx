import { useEffect, useState } from 'react';
import { api, extractList } from '../services/api';

const EMPTY_FORM = {
  title: '',
  destination: '',
  start_date: '',
  end_date: '',
  description: '',
};

export default function ItinerarySection({
  selectedItineraryId,
  onSelectItinerary,
  onNotify,
  refreshKey,
}) {
  const [items, setItems] = useState([]);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(false);

  async function fetchItineraries() {
    setLoading(true);
    try {
      const res = await api.listItineraries();
      const list = extractList(res);
      setItems(list);
      if (!selectedItineraryId && list.length > 0) {
        onSelectItinerary(list[0].id);
      }
    } catch (err) {
      onNotify('error', `Itineraries: ${err.message}`);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchItineraries();
  }, [refreshKey]);

  function startEdit(item) {
    setEditingId(item.id);
    setForm({
      title: item.title || '',
      destination: item.destination || '',
      start_date: item.start_date ? String(item.start_date).slice(0, 10) : '',
      end_date: item.end_date ? String(item.end_date).slice(0, 10) : '',
      description: item.description || '',
    });
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setLoading(true);
    try {
      const payload = {
        title: form.title.trim(),
        destination: form.destination.trim(),
        start_date: form.start_date || undefined,
        end_date: form.end_date || undefined,
        description: form.description.trim(),
      };

      if (editingId) {
        await api.updateItinerary(editingId, payload);
        onNotify('success', 'Itinerary berhasil diperbarui.');
      } else {
        const res = await api.createItinerary(payload);
        const createdId = res?.data?.id || res?.id;
        if (createdId) onSelectItinerary(createdId);
        onNotify('success', 'Itinerary baru berhasil dibuat.');
      }
      cancelEdit();
      await fetchItineraries();
    } catch (err) {
      onNotify('error', err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id) {
    setLoading(true);
    try {
      await api.deleteItinerary(id);
      if (String(selectedItineraryId) === String(id)) {
        onSelectItinerary('');
      }
      onNotify('success', 'Itinerary dihapus.');
      await fetchItineraries();
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
          <h2>2. Itineraries</h2>
          <p className="muted">Pilih rencana perjalanan untuk melihat dan menambah Spots.</p>
        </div>
        <button type="button" className="btn btn-ghost btn-sm" onClick={fetchItineraries}>
          Refresh
        </button>
      </div>

      <form onSubmit={handleSubmit} className="form-grid">
        <label>
          Judul Itinerary *
          <input
            type="text"
            placeholder="Contoh: Liburan ke Bali 3 Hari"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            required
          />
        </label>
        <label>
          Destinasi Utama
          <input
            type="text"
            placeholder="Contoh: Bali"
            value={form.destination}
            onChange={(e) => setForm({ ...form, destination: e.target.value })}
          />
        </label>
        <label>
          Tanggal Mulai
          <input
            type="date"
            value={form.start_date}
            onChange={(e) => setForm({ ...form, start_date: e.target.value })}
          />
        </label>
        <label>
          Tanggal Selesai
          <input
            type="date"
            value={form.end_date}
            onChange={(e) => setForm({ ...form, end_date: e.target.value })}
          />
        </label>
        <label className="span-2">
          Deskripsi
          <input
            type="text"
            placeholder="Catatan singkat rencana perjalanan..."
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
          />
        </label>
        <div className="form-actions span-2">
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {editingId ? 'Simpan Perubahan' : 'Tambah Itinerary'}
          </button>
          {editingId && (
            <button type="button" className="btn btn-ghost" onClick={cancelEdit}>
              Batal
            </button>
          )}
        </div>
      </form>

      <div className="item-list">
        {items.length === 0 ? (
          <p className="empty-state">Belum ada data itinerary.</p>
        ) : (
          items.map((item) => {
            const isSelected = String(item.id) === String(selectedItineraryId);
            return (
              <div key={item.id} className={isSelected ? 'item-row item-selected' : 'item-row'}>
                <div onClick={() => onSelectItinerary(item.id)} className="clickable">
                  <strong>#{item.id} — {item.title}</strong>
                  <div className="muted">
                    {item.destination ? `Tujuan: ${item.destination}` : 'Tanpa tujuan'}
                    {item.start_date ? ` • ${String(item.start_date).slice(0, 10)}` : ''}
                  </div>
                </div>
                <div className="row-actions">
                  <button
                    type="button"
                    className="btn btn-secondary btn-sm"
                    onClick={() => onSelectItinerary(item.id)}
                  >
                    {isSelected ? 'Dipilih' : 'Pilih'}
                  </button>
                  <button
                    type="button"
                    className="btn btn-ghost btn-sm"
                    onClick={() => startEdit(item)}
                  >
                    Edit
                  </button>
                  <button
                    type="button"
                    className="btn btn-danger btn-sm"
                    onClick={() => handleDelete(item.id)}
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
