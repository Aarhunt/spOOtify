import { create } from 'zustand'
import axios from "axios";

const initialState = {
  loading: false,
  success: false,
  error: false,
  data: null,
  errorData: null,
};

const usePlaylistStore = create((set, get) => ({
    ...initialState,

    execute: async () => {
        set({...initialState, loading: true});
        try {
            const res = await axios.get("localhost:8080/playlists");
            set({...initialState, success: true, data: res.data})
        }
    }
}));
