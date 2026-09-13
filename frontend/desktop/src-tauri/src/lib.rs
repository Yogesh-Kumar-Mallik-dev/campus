/**
 * BLOCK_DESKTOP_TAURI_LIB_001
 * Subsystem: The Hub Root Super-App Shell (Desktop Wrapper)
 * Purpose:   Tauri 2 desktop lifecycle and native command bindings.
 */

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}
