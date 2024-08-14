import { createSlice } from "@reduxjs/toolkit";
import type { PayloadAction } from "@reduxjs/toolkit";

const initialState = {
  selectedKey: "",
};

export const appSlice = createSlice({
  name: "app",
  initialState,
  reducers: {
    setRoute: (state, action: PayloadAction<string>): void => {
      state.selectedKey = action.payload;
    },
  },
});

export const { setRoute } = appSlice.actions;

export default appSlice.reducer;
