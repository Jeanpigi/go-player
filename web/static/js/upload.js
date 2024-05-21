document
  .getElementById("uploadForm")
  .addEventListener("submit", function (event) {
    const fileInput = document.getElementById("musicFile");
    const file = fileInput.files[0];

    if (file && file.type !== "audio/mp3") {
      event.preventDefault();
      Swal.fire({
        icon: "error",
        title: "Invalid File",
        text: "Please upload a valid MP3 file.",
      });
    }
  });
