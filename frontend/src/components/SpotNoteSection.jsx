import { useEffect, useState } from 'react';
import { api, extractList } from '../services/api';

const EMPTY_FORM = {
  title: '',
  note_content: '',
  opening_hours: '',
};

export default function SpotNoteSection({ spotId, onNotify, refreshKey }) {
  const [notes, setNotes] = useState([]);
  const [form, setForm] = useState(EMPTY_FORM);
  const [editingId, setEditingId] = useState(null);
  const [loading, setLoading] = useState(false);

  async function fetchNotes() {
    setLoading(true);
    try {
      const res = await api.listSpotNotes(spotId);
      const all = extractList(res);
      const filtered = spotId
        ? all.filter((n) => String(n.spot_id) === String(spotId))
        : all;
      setNotes(filtered);
    } catch (err) {
      onNotify('error', `Spot Notes: ${err.message}`);
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchNotes();
  }, [spotId, refreshKey]);

  function startEdit(note) {
    setEditingId(note.id);
    setForm({
      title: note.title || '',
      note_content: note.note_content || '',
      opening_hours: note.opening_hours || '',
    });
  }

  function cancelEdit() {
    setEditingId(null);
    setForm(EMPTY_FORM);
  }

  async function handleSubmit(e) {
    e.preventDefault();
    if (!spotId) {
      onNotify('error', 'Pilih atau buat Spot terlebih dahulu.');
      return;
    }
    setLoading(true);
    try {
      const payload = {
        spot_id: Number(spotId),
        title: form.title.trim(),
        note_content: form.note_content.trim(),
        opening_hours: form.opening_hours.trim(),
      };
      if (editingId) {
        await api.updateSpotNote(editingId, payload);
        onNotify('success', 'Catatan spot diperbarui.');
      } else {
        await api.createSpotNote(payload);
        onNotify('success', 'Catatan spot disimpan.');
      }
      cancelEdit();
      await fetchNotes();
    } catch (err) {
      onNotify('error', err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleDelete(id) {
    setLoading(true);
    try {
      await api.deleteSpotNote(id);
      onNotify('success', 'Catatan spot dihapus.');
      await fetchNotes();
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
          <h2>4. Spot Notes</h2>
          <p className="muted">
            Spot aktif: <strong>{spotId ? `#${spotId}` : 'Belum dipilih'}</strong>
          </p>
        </div>
        <button type="button" className="btn btn-ghost btn-sm" onClick={fetchNotes}>
          Refresh
        </button>
      </div>

      <form onSubmit={handleSubmit} className="form-grid">
        <label>
          Judul Catatan *
          <input
            type="text"
            placeholder="Contoh: Tips Datang Pagi"
            value={form.title}
            onChange={(e) => setForm({ ...form, title: e.target.value })}
            required
          />
        </label>
        <label>
          Jam Operasional
          <input
            type="text"
            placeholder="Contoh: 07.00 - 17.00"
            value={form.opening_hours}
            onChange={(e) => setForm({ ...form, opening_hours: e.target.value })}
          />
        </label>
        <label className="span-2">
          Isi Catatan *
          <textarea
            rows="3"
            placeholder="Tulis tips kunjungan, tiket, atau hal penting..."
            value={form.note_content}
            onChange={(e) => setForm({ ...form, note_content: e.target.value })}
            required
          />
        </label>
        <div className="form-actions span-2">
          <button type="submit" className="btn btn-primary" disabled={loading || !spotId}>
            {editingId ? 'Simpan Catatan' : 'Tambah Catatan'}
          </button>
          {editingId && (
            <button type="button" className="btn btn-ghost" onClick={cancelEdit}>
              Batal
            </button>
          )}
        </div>
      </form>

      <div className="item-list">
        {notes.length === 0 ? (
          <p className="empty-state">Belum ada catatan untuk spot ini.</p>
        ) : (
          notes.map((note) => (
            <div key={note.id} className="item-row">
              <div>
                <strong>#{note.id} — {note.title}</strong>
                {note.opening_hours && (
                  <span className="badge">Jam: {note.opening_hours}</span>
                )}
                <p className="note-body">{note.note_content}</p>
                {note.spot?.name && (
                  <div className="muted">Relasi Spot: {note.spot.name}</div>
                )}
              </div>
              <div className="row-actions">
                <button
                  type="button"
                  className="btn btn-ghost btn-sm"
                  onClick={() => startEdit(note)}
                >
                  Edit
                </button>
                <button
                  type="button"
                  className="btn btn-danger btn-sm"
                  onClick={() => handleDelete(note.id)}
                >
                  Hapus
                </button>
              </div>
            </div>
          ))
        )}
      </div>
    </section>
  );
}
