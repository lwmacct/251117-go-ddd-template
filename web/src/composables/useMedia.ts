/**
 * Media Composable
 * 提供媒体（音频/视频）相关的工具函数
 */

import {
  ref,
  computed,
  watch,
  onMounted,
  onUnmounted,
  type Ref,
  type ComputedRef,
} from "vue";

// ============================================================================
// 类型定义
// ============================================================================

/**
 * 媒体状态
 */
export type MediaState = "idle" | "loading" | "ready" | "playing" | "paused" | "ended" | "error";

/**
 * 媒体选项
 */
export interface UseMediaOptions {
  /** 是否自动播放 */
  autoplay?: boolean;
  /** 是否循环 */
  loop?: boolean;
  /** 是否静音 */
  muted?: boolean;
  /** 初始音量 (0-1) */
  volume?: number;
  /** 播放速度 */
  playbackRate?: number;
  /** 是否预加载 */
  preload?: "none" | "metadata" | "auto";
}

/**
 * 媒体返回值
 */
export interface UseMediaReturn {
  /** 媒体元素引用 */
  mediaRef: Ref<HTMLMediaElement | null>;
  /** 当前状态 */
  state: Ref<MediaState>;
  /** 是否正在播放 */
  isPlaying: ComputedRef<boolean>;
  /** 是否暂停 */
  isPaused: ComputedRef<boolean>;
  /** 是否静音 */
  isMuted: Ref<boolean>;
  /** 当前时间（秒） */
  currentTime: Ref<number>;
  /** 总时长（秒） */
  duration: Ref<number>;
  /** 缓冲进度 (0-1) */
  buffered: Ref<number>;
  /** 播放进度 (0-1) */
  progress: ComputedRef<number>;
  /** 音量 (0-1) */
  volume: Ref<number>;
  /** 播放速度 */
  playbackRate: Ref<number>;
  /** 是否可播放 */
  canPlay: Ref<boolean>;
  /** 是否在等待数据 */
  isWaiting: Ref<boolean>;
  /** 是否在缓冲 */
  isBuffering: ComputedRef<boolean>;
  /** 错误信息 */
  error: Ref<MediaError | null>;
  /** 播放 */
  play: () => Promise<void>;
  /** 暂停 */
  pause: () => void;
  /** 停止 */
  stop: () => void;
  /** 跳转 */
  seek: (time: number) => void;
  /** 跳转百分比 */
  seekPercent: (percent: number) => void;
  /** 切换播放/暂停 */
  toggle: () => Promise<void>;
  /** 切换静音 */
  toggleMute: () => void;
  /** 快进 */
  forward: (seconds?: number) => void;
  /** 快退 */
  backward: (seconds?: number) => void;
  /** 设置源 */
  setSource: (src: string) => void;
}

/**
 * 音频可视化选项
 */
export interface AudioVisualizerOptions {
  /** FFT 大小 */
  fftSize?: number;
  /** 平滑系数 */
  smoothingTimeConstant?: number;
}

/**
 * 音频可视化返回值
 */
export interface AudioVisualizerReturn {
  /** 连接到音频元素 */
  connect: (audio: HTMLAudioElement) => void;
  /** 断开连接 */
  disconnect: () => void;
  /** 频率数据 */
  frequencyData: Ref<Uint8Array>;
  /** 时域数据 */
  timeDomainData: Ref<Uint8Array>;
  /** 是否活跃 */
  isActive: Ref<boolean>;
}

/**
 * 录音选项
 */
export interface RecorderOptions {
  /** MIME 类型 */
  mimeType?: string;
  /** 音频比特率 */
  audioBitsPerSecond?: number;
  /** 视频比特率 */
  videoBitsPerSecond?: number;
}

/**
 * 录音返回值
 */
export interface UseRecorderReturn {
  /** 是否正在录制 */
  isRecording: Ref<boolean>;
  /** 是否暂停 */
  isPaused: Ref<boolean>;
  /** 录制时长（秒） */
  duration: Ref<number>;
  /** 录制数据 */
  data: Ref<Blob | null>;
  /** 数据 URL */
  dataUrl: Ref<string | null>;
  /** 开始录制 */
  start: () => Promise<void>;
  /** 暂停录制 */
  pause: () => void;
  /** 恢复录制 */
  resume: () => void;
  /** 停止录制 */
  stop: () => Promise<Blob>;
  /** 清除数据 */
  clear: () => void;
  /** 错误 */
  error: Ref<Error | null>;
}

// ============================================================================
// 核心函数
// ============================================================================

/**
 * 使用媒体
 *
 * @description 创建媒体控制器
 *
 * @example
 * ```ts
 * const {
 *   mediaRef,
 *   play,
 *   pause,
 *   toggle,
 *   isPlaying,
 *   currentTime,
 *   duration,
 *   progress,
 *   volume,
 *   seek
 * } = useMedia({
 *   autoplay: false,
 *   volume: 0.8
 * })
 *
 * // 在模板中: <video ref="mediaRef" :src="videoUrl" />
 * ```
 */
export function useMedia(options: UseMediaOptions = {}): UseMediaReturn {
  const {
    autoplay = false,
    loop = false,
    muted = false,
    volume: initialVolume = 1,
    playbackRate: initialPlaybackRate = 1,
    preload = "metadata",
  } = options;

  const mediaRef = ref<HTMLMediaElement | null>(null);
  const state = ref<MediaState>("idle");
  const currentTime = ref(0);
  const duration = ref(0);
  const buffered = ref(0);
  const volume = ref(initialVolume);
  const isMuted = ref(muted);
  const playbackRate = ref(initialPlaybackRate);
  const canPlay = ref(false);
  const isWaiting = ref(false);
  const error = ref<MediaError | null>(null);

  const isPlaying = computed(() => state.value === "playing");
  const isPaused = computed(() => state.value === "paused");
  const progress = computed(() => (duration.value > 0 ? currentTime.value / duration.value : 0));
  const isBuffering = computed(() => isWaiting.value && isPlaying.value);

  // 事件处理器
  const handleLoadStart = () => {
    state.value = "loading";
  };

  const handleCanPlay = () => {
    canPlay.value = true;
    if (state.value === "loading") {
      state.value = "ready";
    }
  };

  const handlePlay = () => {
    state.value = "playing";
  };

  const handlePause = () => {
    if (state.value !== "ended") {
      state.value = "paused";
    }
  };

  const handleEnded = () => {
    state.value = "ended";
  };

  const handleTimeUpdate = () => {
    if (mediaRef.value) {
      currentTime.value = mediaRef.value.currentTime;
    }
  };

  const handleDurationChange = () => {
    if (mediaRef.value) {
      duration.value = mediaRef.value.duration;
    }
  };

  const handleProgress = () => {
    if (mediaRef.value && mediaRef.value.buffered.length > 0) {
      const bufferedEnd = mediaRef.value.buffered.end(mediaRef.value.buffered.length - 1);
      buffered.value = duration.value > 0 ? bufferedEnd / duration.value : 0;
    }
  };

  const handleVolumeChange = () => {
    if (mediaRef.value) {
      volume.value = mediaRef.value.volume;
      isMuted.value = mediaRef.value.muted;
    }
  };

  const handleWaiting = () => {
    isWaiting.value = true;
  };

  const handlePlaying = () => {
    isWaiting.value = false;
  };

  const handleError = () => {
    state.value = "error";
    if (mediaRef.value) {
      error.value = mediaRef.value.error;
    }
  };

  // 设置媒体元素
  const setupMedia = (media: HTMLMediaElement) => {
    media.autoplay = autoplay;
    media.loop = loop;
    media.muted = muted;
    media.volume = initialVolume;
    media.playbackRate = initialPlaybackRate;
    media.preload = preload;

    media.addEventListener("loadstart", handleLoadStart);
    media.addEventListener("canplay", handleCanPlay);
    media.addEventListener("play", handlePlay);
    media.addEventListener("pause", handlePause);
    media.addEventListener("ended", handleEnded);
    media.addEventListener("timeupdate", handleTimeUpdate);
    media.addEventListener("durationchange", handleDurationChange);
    media.addEventListener("progress", handleProgress);
    media.addEventListener("volumechange", handleVolumeChange);
    media.addEventListener("waiting", handleWaiting);
    media.addEventListener("playing", handlePlaying);
    media.addEventListener("error", handleError);
  };

  const cleanupMedia = (media: HTMLMediaElement) => {
    media.removeEventListener("loadstart", handleLoadStart);
    media.removeEventListener("canplay", handleCanPlay);
    media.removeEventListener("play", handlePlay);
    media.removeEventListener("pause", handlePause);
    media.removeEventListener("ended", handleEnded);
    media.removeEventListener("timeupdate", handleTimeUpdate);
    media.removeEventListener("durationchange", handleDurationChange);
    media.removeEventListener("progress", handleProgress);
    media.removeEventListener("volumechange", handleVolumeChange);
    media.removeEventListener("waiting", handleWaiting);
    media.removeEventListener("playing", handlePlaying);
    media.removeEventListener("error", handleError);
  };

  // 监听 mediaRef 变化
  watch(mediaRef, (newMedia, oldMedia) => {
    if (oldMedia) {
      cleanupMedia(oldMedia);
    }
    if (newMedia) {
      setupMedia(newMedia);
    }
  });

  // 同步属性到媒体元素
  watch(volume, (v) => {
    if (mediaRef.value) {
      mediaRef.value.volume = v;
    }
  });

  watch(isMuted, (v) => {
    if (mediaRef.value) {
      mediaRef.value.muted = v;
    }
  });

  watch(playbackRate, (v) => {
    if (mediaRef.value) {
      mediaRef.value.playbackRate = v;
    }
  });

  // 控制方法
  const play = async (): Promise<void> => {
    if (mediaRef.value) {
      await mediaRef.value.play();
    }
  };

  const pause = () => {
    if (mediaRef.value) {
      mediaRef.value.pause();
    }
  };

  const stop = () => {
    if (mediaRef.value) {
      mediaRef.value.pause();
      mediaRef.value.currentTime = 0;
      state.value = "ready";
    }
  };

  const seek = (time: number) => {
    if (mediaRef.value) {
      mediaRef.value.currentTime = Math.max(0, Math.min(time, duration.value));
    }
  };

  const seekPercent = (percent: number) => {
    seek(duration.value * Math.max(0, Math.min(1, percent)));
  };

  const toggle = async (): Promise<void> => {
    if (isPlaying.value) {
      pause();
    } else {
      await play();
    }
  };

  const toggleMute = () => {
    isMuted.value = !isMuted.value;
  };

  const forward = (seconds = 10) => {
    seek(currentTime.value + seconds);
  };

  const backward = (seconds = 10) => {
    seek(currentTime.value - seconds);
  };

  const setSource = (src: string) => {
    if (mediaRef.value) {
      mediaRef.value.src = src;
      mediaRef.value.load();
      state.value = "idle";
    }
  };

  onUnmounted(() => {
    if (mediaRef.value) {
      cleanupMedia(mediaRef.value);
    }
  });

  return {
    mediaRef,
    state,
    isPlaying,
    isPaused,
    isMuted,
    currentTime,
    duration,
    buffered,
    progress,
    volume,
    playbackRate,
    canPlay,
    isWaiting,
    isBuffering,
    error,
    play,
    pause,
    stop,
    seek,
    seekPercent,
    toggle,
    toggleMute,
    forward,
    backward,
    setSource,
  };
}

/**
 * 使用音频
 *
 * @description 创建音频播放器
 *
 * @example
 * ```ts
 * const { play, pause, isPlaying, currentTime, duration } = useAudio('/audio.mp3', {
 *   autoplay: false
 * })
 *
 * play()
 * ```
 */
export function useAudio(
  src?: string,
  options: UseMediaOptions = {}
): UseMediaReturn & { audio: Ref<HTMLAudioElement | null> } {
  const media = useMedia(options);
  const audio = ref<HTMLAudioElement | null>(null);

  onMounted(() => {
    const audioElement = new Audio();
    if (src) {
      audioElement.src = src;
    }
    audio.value = audioElement;
    media.mediaRef.value = audioElement;
  });

  onUnmounted(() => {
    if (audio.value) {
      audio.value.pause();
      audio.value.src = "";
      audio.value = null;
    }
  });

  return {
    ...media,
    audio,
  };
}

/**
 * 使用音频可视化
 *
 * @description 创建音频可视化分析器
 *
 * @example
 * ```ts
 * const { connect, frequencyData, timeDomainData, isActive } = useAudioVisualizer({
 *   fftSize: 256
 * })
 *
 * // 连接到音频元素
 * const audioEl = document.querySelector('audio')
 * connect(audioEl)
 *
 * // 使用频率数据绘制可视化
 * console.log(frequencyData.value) // Uint8Array
 * ```
 */
export function useAudioVisualizer(
  options: AudioVisualizerOptions = {}
): AudioVisualizerReturn {
  const { fftSize = 256, smoothingTimeConstant = 0.8 } = options;

  const frequencyData = ref<Uint8Array>(new Uint8Array(fftSize / 2));
  const timeDomainData = ref<Uint8Array>(new Uint8Array(fftSize / 2));
  const isActive = ref(false);

  let audioContext: AudioContext | null = null;
  let analyser: AnalyserNode | null = null;
  let source: MediaElementAudioSourceNode | null = null;
  let animationFrame: number | null = null;

  const updateData = () => {
    if (!analyser || !isActive.value) return;

    analyser.getByteFrequencyData(frequencyData.value);
    analyser.getByteTimeDomainData(timeDomainData.value);

    animationFrame = requestAnimationFrame(updateData);
  };

  const connect = (audio: HTMLAudioElement) => {
    if (audioContext) {
      disconnect();
    }

    audioContext = new AudioContext();
    analyser = audioContext.createAnalyser();
    analyser.fftSize = fftSize;
    analyser.smoothingTimeConstant = smoothingTimeConstant;

    source = audioContext.createMediaElementSource(audio);
    source.connect(analyser);
    analyser.connect(audioContext.destination);

    frequencyData.value = new Uint8Array(analyser.frequencyBinCount);
    timeDomainData.value = new Uint8Array(analyser.frequencyBinCount);

    isActive.value = true;
    updateData();
  };

  const disconnect = () => {
    isActive.value = false;

    if (animationFrame) {
      cancelAnimationFrame(animationFrame);
      animationFrame = null;
    }

    if (source) {
      source.disconnect();
      source = null;
    }

    if (analyser) {
      analyser.disconnect();
      analyser = null;
    }

    if (audioContext) {
      audioContext.close();
      audioContext = null;
    }
  };

  onUnmounted(disconnect);

  return {
    connect,
    disconnect,
    frequencyData,
    timeDomainData,
    isActive,
  };
}

/**
 * 使用录音机
 *
 * @description 创建音频/视频录制器
 *
 * @example
 * ```ts
 * const { start, stop, pause, resume, isRecording, data, dataUrl } = useRecorder({
 *   mimeType: 'audio/webm'
 * })
 *
 * // 开始录音
 * await start()
 *
 * // 停止并获取数据
 * const blob = await stop()
 *
 * // 播放录音
 * const audio = new Audio(dataUrl.value)
 * audio.play()
 * ```
 */
export function useRecorder(options: RecorderOptions = {}): UseRecorderReturn {
  const {
    mimeType = "audio/webm",
    audioBitsPerSecond = 128000,
    videoBitsPerSecond,
  } = options;

  const isRecording = ref(false);
  const isPaused = ref(false);
  const duration = ref(0);
  const data = ref<Blob | null>(null);
  const dataUrl = ref<string | null>(null);
  const error = ref<Error | null>(null);

  let mediaRecorder: MediaRecorder | null = null;
  let chunks: Blob[] = [];
  let stream: MediaStream | null = null;
  let startTime = 0;
  let durationInterval: ReturnType<typeof setInterval> | null = null;

  const start = async (): Promise<void> => {
    try {
      error.value = null;
      chunks = [];

      // 获取媒体流
      const isAudio = mimeType.startsWith("audio/");
      stream = await navigator.mediaDevices.getUserMedia({
        audio: true,
        video: !isAudio,
      });

      // 创建录制器
      const recorderOptions: MediaRecorderOptions = {
        mimeType,
        audioBitsPerSecond,
      };

      if (videoBitsPerSecond) {
        recorderOptions.videoBitsPerSecond = videoBitsPerSecond;
      }

      mediaRecorder = new MediaRecorder(stream, recorderOptions);

      mediaRecorder.ondataavailable = (e) => {
        if (e.data.size > 0) {
          chunks.push(e.data);
        }
      };

      mediaRecorder.start();
      isRecording.value = true;
      isPaused.value = false;
      startTime = Date.now();

      // 更新时长
      durationInterval = setInterval(() => {
        if (!isPaused.value) {
          duration.value = (Date.now() - startTime) / 1000;
        }
      }, 100);
    } catch (e) {
      error.value = e instanceof Error ? e : new Error(String(e));
      throw error.value;
    }
  };

  const pause = () => {
    if (mediaRecorder && mediaRecorder.state === "recording") {
      mediaRecorder.pause();
      isPaused.value = true;
    }
  };

  const resume = () => {
    if (mediaRecorder && mediaRecorder.state === "paused") {
      mediaRecorder.resume();
      isPaused.value = false;
    }
  };

  const stop = async (): Promise<Blob> => {
    return new Promise((resolve, reject) => {
      if (!mediaRecorder) {
        reject(new Error("Recorder not initialized"));
        return;
      }

      mediaRecorder.onstop = () => {
        // 创建 Blob
        const blob = new Blob(chunks, { type: mimeType });
        data.value = blob;

        // 创建 URL
        if (dataUrl.value) {
          URL.revokeObjectURL(dataUrl.value);
        }
        dataUrl.value = URL.createObjectURL(blob);

        // 停止流
        if (stream) {
          stream.getTracks().forEach((track) => track.stop());
          stream = null;
        }

        // 清理
        if (durationInterval) {
          clearInterval(durationInterval);
          durationInterval = null;
        }

        isRecording.value = false;
        isPaused.value = false;
        mediaRecorder = null;

        resolve(blob);
      };

      mediaRecorder.stop();
    });
  };

  const clear = () => {
    data.value = null;
    if (dataUrl.value) {
      URL.revokeObjectURL(dataUrl.value);
      dataUrl.value = null;
    }
    duration.value = 0;
    chunks = [];
  };

  onUnmounted(() => {
    if (mediaRecorder && mediaRecorder.state !== "inactive") {
      mediaRecorder.stop();
    }
    if (stream) {
      stream.getTracks().forEach((track) => track.stop());
    }
    if (durationInterval) {
      clearInterval(durationInterval);
    }
    if (dataUrl.value) {
      URL.revokeObjectURL(dataUrl.value);
    }
  });

  return {
    isRecording,
    isPaused,
    duration,
    data,
    dataUrl,
    start,
    pause,
    resume,
    stop,
    clear,
    error,
  };
}

/**
 * 使用屏幕共享
 *
 * @description 创建屏幕共享捕获
 *
 * @example
 * ```ts
 * const { start, stop, stream, isSharing, error } = useScreenShare()
 *
 * // 开始共享
 * await start({ video: true, audio: true })
 *
 * // 获取视频流
 * video.srcObject = stream.value
 * ```
 */
export function useScreenShare(): {
  stream: Ref<MediaStream | null>;
  isSharing: Ref<boolean>;
  error: Ref<Error | null>;
  start: (options?: DisplayMediaStreamOptions) => Promise<MediaStream>;
  stop: () => void;
} {
  const stream = ref<MediaStream | null>(null);
  const isSharing = ref(false);
  const error = ref<Error | null>(null);

  const start = async (
    options: DisplayMediaStreamOptions = { video: true, audio: false }
  ): Promise<MediaStream> => {
    try {
      error.value = null;
      const mediaStream = await navigator.mediaDevices.getDisplayMedia(options);

      // 监听流结束
      mediaStream.getVideoTracks()[0].onended = () => {
        stop();
      };

      stream.value = mediaStream;
      isSharing.value = true;

      return mediaStream;
    } catch (e) {
      error.value = e instanceof Error ? e : new Error(String(e));
      throw error.value;
    }
  };

  const stop = () => {
    if (stream.value) {
      stream.value.getTracks().forEach((track) => track.stop());
      stream.value = null;
    }
    isSharing.value = false;
  };

  onUnmounted(stop);

  return {
    stream,
    isSharing,
    error,
    start,
    stop,
  };
}

/**
 * 格式化时长
 *
 * @description 将秒数格式化为时间字符串
 *
 * @example
 * ```ts
 * formatDuration(125) // '2:05'
 * formatDuration(3725) // '1:02:05'
 * ```
 */
export function formatDuration(seconds: number): string {
  const h = Math.floor(seconds / 3600);
  const m = Math.floor((seconds % 3600) / 60);
  const s = Math.floor(seconds % 60);

  if (h > 0) {
    return `${h}:${m.toString().padStart(2, "0")}:${s.toString().padStart(2, "0")}`;
  }
  return `${m}:${s.toString().padStart(2, "0")}`;
}

/**
 * 使用画中画
 *
 * @description 控制视频画中画模式
 *
 * @example
 * ```ts
 * const { enter, exit, toggle, isSupported, isPip } = usePictureInPicture(videoRef)
 *
 * if (isSupported.value) {
 *   toggle()
 * }
 * ```
 */
export function usePictureInPicture(videoRef: Ref<HTMLVideoElement | null>): {
  isSupported: ComputedRef<boolean>;
  isPip: Ref<boolean>;
  enter: () => Promise<void>;
  exit: () => Promise<void>;
  toggle: () => Promise<void>;
} {
  const isPip = ref(false);

  const isSupported = computed(() => {
    return "pictureInPictureEnabled" in document && document.pictureInPictureEnabled;
  });

  const enter = async () => {
    if (videoRef.value && isSupported.value) {
      await videoRef.value.requestPictureInPicture();
      isPip.value = true;
    }
  };

  const exit = async () => {
    if (document.pictureInPictureElement) {
      await document.exitPictureInPicture();
      isPip.value = false;
    }
  };

  const toggle = async () => {
    if (isPip.value) {
      await exit();
    } else {
      await enter();
    }
  };

  // 监听画中画状态变化
  const handleEnter = () => {
    isPip.value = true;
  };

  const handleLeave = () => {
    isPip.value = false;
  };

  watch(
    videoRef,
    (video, oldVideo) => {
      if (oldVideo) {
        oldVideo.removeEventListener("enterpictureinpicture", handleEnter);
        oldVideo.removeEventListener("leavepictureinpicture", handleLeave);
      }
      if (video) {
        video.addEventListener("enterpictureinpicture", handleEnter);
        video.addEventListener("leavepictureinpicture", handleLeave);
      }
    },
    { immediate: true }
  );

  return {
    isSupported,
    isPip,
    enter,
    exit,
    toggle,
  };
}
