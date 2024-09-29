document.addEventListener("DOMContentLoaded", () => {
  const video = document.getElementById("video");
  const videoSrc = "http://localhost:8080/stream/stream.m3u8";

  if (Hls.isSupported()) {
    const hls = new Hls();
    hls.loadSource(videoSrc);
    hls.attachMedia(video);
    hls.on(Hls.Events.MANIFEST_PARSED, () => {
      video.play();
    });

    hls.on(Hls.Events.ERROR, (event, data) => {
      const errorType = data.type;
      const errorDetails = data.details;
      const errorFatal = data.fatal;

      console.error("Error HLS:", errorType, errorDetails);

      if (errorFatal) {
        switch (data.type) {
          case Hls.ErrorTypes.NETWORK_ERROR:
            console.error("Error NETWORK_ERROR - try to reload stream");
            hls.startLoad();
            break;
          case Hls.ErrorTypes.MEDIA_ERROR:
            console.error("Error MEDIA_ERROR");
            hls.recoverMediaError();
            break;
          default:
            hls.destroy();
            break;
        }
      }
    });
  } else if (video.canPlayType("application/vnd.apple.mpegurl")) {
    video.src = videoSrc;
    video.addEventListener("loadedmetadata", () => {
      video.play();
    });
  } else {
    alert("Your browser does not support HLS.");
  }
});
