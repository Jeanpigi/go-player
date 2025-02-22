document.getElementById("uploadForm").addEventListener("submit", function (event) {
  const fileInput = document.getElementById("musicFile");
  const file = fileInput.files[0];

  const submitButton = document.getElementById("submitButton");
  // Deshabilitar el botón para evitar múltiples envíos
  submitButton.disabled = true;
  // Cambiar el texto del botón a "Cargando..."
  submitButton.innerText = "Cargando...";

  // Ajusta la verificación al MIME type correcto para MP3
  if (file && file.type !== "audio/mpeg") {
    event.preventDefault();
    // Reactivar el botón y restaurar el texto en caso de error
    submitButton.disabled = false;
    submitButton.innerText = "Upload";
    Swal.fire({
      icon: "error",
      title: "Invalid File",
      text: "Please upload a valid MP3 file.",
    });
  }
});


