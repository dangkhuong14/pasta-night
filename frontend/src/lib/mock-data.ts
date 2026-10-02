// Mock API data, used when NEXT_PUBLIC_API_BASE_URL is unset (frontend/CLAUDE.md).
// Movies are copied from a real refresh of the backend cache (option "friends").
// Some fields are altered on purpose so every hide rule and fallback in
// docs/SCREENS.md shows up in mock mode:
// - 969681: last cast member has no photo (initials fallback)
// - 1368337: no poster, no backdrop
// - 1204680: no runtime, no release year
// - 1032863: no providers (the "CÓ MẶT TRÊN" card is hidden)
// - 1288445: empty cast
// - 1032863: media has no trailer (the carousel opens on a still)
// - 1368337: media is empty (the hero falls back to the poster placeholder)
// Carousel stills are real TMDB backdrops, but each movie borrows two from
// its neighbours so there is more than one slide to scroll through.
// Titles with an empty overview are real TMDB gaps (no vi-VN translation).
import type { MovieDetail, ViewingOption } from "./api-types";

export const MOCK_FETCHED_AT = "2026-09-29T11:35:40Z";

export const MOCK_OPTIONS: ViewingOption[] = [
  {
    id: "netflix-chill",
    label: "Netflix & Chill",
    description: "Cho hai người",
    icon: "heart",
  },
  {
    id: "solo",
    label: "Một mình",
    description: "Thời gian cho riêng bạn",
    icon: "user",
  },
  {
    id: "friends",
    label: "Hội bạn",
    description: "Xem cùng nhóm bạn",
    icon: "users",
  },
];

export const MOCK_MOVIES: MovieDetail[] = [
  {
    id: 969681,
    title: "Người Nhện: Khởi Đầu Mới",
    original_title: "Spider-Man: Brand New Day",
    overview:
      "Không còn Tony Stark, MJ hay Ned kề cận, Peter buộc phải đơn thân độc mã đối diện với phe đối đầu bí ẩn. Tuy nhiên, khi áp lực ngày càng gia tăng, nó kích hoạt một sự biến đổi thể chất bất ngờ, đe dọa chính sự tồn tại của anh. Đồng thời, một chuỗi tội phạm bí ẩn mới xuất hiện, kéo theo một trong những mối đe dọa mạnh mẽ nhất mà Spider-Man từng đối mặt.",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/wqGZVSCUSXE92WH2zyol2REaqT4.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/qeQJx07rK2xm8SD2sJxFKhE7gs0.jpg",
    release_year: 2026,
    rating: 7.9,
    vote_count: 2930,
    runtime_minutes: 150,
    genres: ["Phim Khoa Học Viễn Tưởng", "Phim Hành Động", "Phim Phiêu Lưu"],
    directors: ["Destin Daniel Cretton"],
    cast: [
      {
        name: "Tom Holland",
        profile_url:
          "https://image.tmdb.org/t/p/w185/xKBAaPIa1c7tzZD3Y0MhBLv4hPE.jpg",
      },
      {
        name: "Zendaya",
        profile_url:
          "https://image.tmdb.org/t/p/w185/1qup8tSt95HLbcy2c2xrx4iJNxv.jpg",
      },
      {
        name: "Mark Ruffalo",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5GilHMOt5PAQh6rlUKZzGmaKEI7.jpg",
      },
      {
        name: "Jon Bernthal",
        profile_url:
          "https://image.tmdb.org/t/p/w185/bjcglF1IpDcZNs1HwFxakGpzBo6.jpg",
      },
      {
        name: "Jacob Batalon",
        profile_url: null,
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=P3uI5sLosKU",
        youtube_key: "P3uI5sLosKU",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/qeQJx07rK2xm8SD2sJxFKhE7gs0.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/3icyRAqgakNcQn6aDVz9libFmBA.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/viZqGq9TNvQ5uXSD4ahg2RpRONT.jpg",
        youtube_key: null,
      },
    ],
    providers: [],
  },
  {
    id: 1423191,
    title: "Vùng Đất Quỷ Dữ",
    original_title: "Resident Evil",
    overview: "",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/jlFMvLuzjwt7aI1QVtiNVnFHNBe.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/3icyRAqgakNcQn6aDVz9libFmBA.jpg",
    release_year: 2026,
    rating: 7.3,
    vote_count: 583,
    runtime_minutes: 95,
    genres: ["Phim Kinh Dị", "Phim Khoa Học Viễn Tưởng", "Phim Phiêu Lưu"],
    directors: ["Zach Cregger"],
    cast: [
      {
        name: "Austin Abrams",
        profile_url:
          "https://image.tmdb.org/t/p/w185/nvbJNwyhIICUFcS4lhxj9mGvL5z.jpg",
      },
      {
        name: "Paul Walter Hauser",
        profile_url:
          "https://image.tmdb.org/t/p/w185/hXjjbYg1Ah8mFf5ZcaakyXzDKMx.jpg",
      },
      {
        name: "Kali Reis",
        profile_url:
          "https://image.tmdb.org/t/p/w185/ruLDXHnKA4aHEQCcFGDAJ7bLTdt.jpg",
      },
      {
        name: "Zach Cherry",
        profile_url:
          "https://image.tmdb.org/t/p/w185/fT3Wv8ef0Vn0daHWAObCp2Bd4Y.jpg",
      },
      {
        name: "Johnno Wilson",
        profile_url:
          "https://image.tmdb.org/t/p/w185/eZtBoBE0F7Qb9ZyNpacMOQtmWnM.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=mNd1gb19A-c",
        youtube_key: "mNd1gb19A-c",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/3icyRAqgakNcQn6aDVz9libFmBA.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/viZqGq9TNvQ5uXSD4ahg2RpRONT.jpg",
        youtube_key: null,
      },
    ],
    providers: [],
  },
  {
    id: 1368337,
    title: "The Odyssey",
    original_title: "The Odyssey",
    overview:
      "Câu chuyện theo chân Odysseus trong hành trình kéo dài 10 năm trở về nhà sau cuộc chiến thành Troy, nơi ông phải đối mặt với các vị thần, quái vật và vô vàn thử thách, đồng thời nỗ lực đoàn tụ với vợ và giành lại vương quốc của mình.",
    tagline: "Thách thức thần linh.",
    poster_url: null,
    backdrop_url: null,
    release_year: 2026,
    rating: 8,
    vote_count: 3964,
    runtime_minutes: 173,
    genres: ["Phim Phiêu Lưu", "Phim Hành Động", "Phim Giả Tượng"],
    directors: ["Christopher Nolan"],
    cast: [
      {
        name: "Matt Damon",
        profile_url:
          "https://image.tmdb.org/t/p/w185/aCvBXTAR9B1qRjIRzMBYhhbm1fR.jpg",
      },
      {
        name: "Tom Holland",
        profile_url:
          "https://image.tmdb.org/t/p/w185/xKBAaPIa1c7tzZD3Y0MhBLv4hPE.jpg",
      },
      {
        name: "Anne Hathaway",
        profile_url:
          "https://image.tmdb.org/t/p/w185/nbccV2pMoyLTCeg5DQip24Eq0Jp.jpg",
      },
      {
        name: "Robert Pattinson",
        profile_url:
          "https://image.tmdb.org/t/p/w185/sRUM2u8qLcsOaTm0jGJGlOEQhlQ.jpg",
      },
      {
        name: "Himesh Patel",
        profile_url:
          "https://image.tmdb.org/t/p/w185/icqsXLmU0FxBGTv63kkA0GcrecO.jpg",
      },
    ],
    media: [],
    providers: [],
  },
  {
    id: 1204680,
    title: "Coyote vs. Acme",
    original_title: "Coyote vs. Acme",
    overview:
      "Sau quá nhiều lần thất bại vì những món hàng lỗi của công ty Acme trong cuộc săn đuổi không ngừng nghỉ chú chim Roadrunner, Wile E. Coyote quyết định thuê một luật sư… từ biển quảng cáo, để kiện tập đoàn Acme ra tòa.",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/kYDCl2y0VPvhT5eYWbMRInPoB03.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/viZqGq9TNvQ5uXSD4ahg2RpRONT.jpg",
    release_year: null,
    rating: 7.6,
    vote_count: 553,
    runtime_minutes: null,
    genres: ["Phim Hài", "Phim Phiêu Lưu", "Phim Gia Đình"],
    directors: ["Dave Green"],
    cast: [
      {
        name: "Will Forte",
        profile_url:
          "https://image.tmdb.org/t/p/w185/4VEzbkL3HwHTUZAPA5PyypFG2U.jpg",
      },
      {
        name: "Trần Đồng Lan",
        profile_url:
          "https://image.tmdb.org/t/p/w185/vWn27Fk2GLwH7o9fBG9hBWZI6OR.jpg",
      },
      {
        name: "John Cena",
        profile_url:
          "https://image.tmdb.org/t/p/w185/rgB2eIOt7WyQjdgJCOuESdDlrjg.jpg",
      },
      {
        name: "Tone Bell",
        profile_url:
          "https://image.tmdb.org/t/p/w185/fB06Xj5FgwvHcZESsNGMupGBTyY.jpg",
      },
      {
        name: "Martha Kelly",
        profile_url:
          "https://image.tmdb.org/t/p/w185/ac7HoARDIHdYNqpzgJBzpmKHXHR.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=kMsiD1Nky5I",
        youtube_key: "kMsiD1Nky5I",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/viZqGq9TNvQ5uXSD4ahg2RpRONT.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/o7Oy9Gbx1CCyaweL8xUhtMW4Puq.jpg",
        youtube_key: null,
      },
    ],
    providers: [],
  },
  {
    id: 1032863,
    title: "The Love Hypothesis",
    original_title: "The Love Hypothesis",
    overview: "",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/vfZxVHextAGC70zrNhS8lsROqP1.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/o7Oy9Gbx1CCyaweL8xUhtMW4Puq.jpg",
    release_year: 2026,
    rating: 8.1,
    vote_count: 353,
    runtime_minutes: 112,
    genres: ["Phim Lãng Mạn", "Phim Hài"],
    directors: ["Claire Scanlon"],
    cast: [
      {
        name: "Lili Reinhart",
        profile_url:
          "https://image.tmdb.org/t/p/w185/sJcuQAJdIv3mnh9M9p5sQWhWSRN.jpg",
      },
      {
        name: "Tom Bateman",
        profile_url:
          "https://image.tmdb.org/t/p/w185/6elQzmY4EJHlLvHsaeoeo2dXOQ1.jpg",
      },
      {
        name: "Rachel Marsh",
        profile_url:
          "https://image.tmdb.org/t/p/w185/6R0pDuQuSf4CEtZHCRnZkqc2EAN.jpg",
      },
      {
        name: "Jaboukie Young-White",
        profile_url:
          "https://image.tmdb.org/t/p/w185/8OYI8OqsMUpBZica4lM2odjgVH6.jpg",
      },
      {
        name: "Nicholas Duvernay",
        profile_url:
          "https://image.tmdb.org/t/p/w185/4uwcWYbXxrxyhJbCPXshTQXQQKi.jpg",
      },
    ],
    media: [
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/o7Oy9Gbx1CCyaweL8xUhtMW4Puq.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/e2QAGrEmbpmZpMymDRkDisJkvg9.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/hpBGCnzOvdtQoMyE48gvwp2y5yx.jpg",
        youtube_key: null,
      },
    ],
    providers: [],
  },
  {
    id: 1288445,
    title: "Tàu Buôn Người",
    original_title: "Mutiny",
    overview:
      "Sau khi chứng kiến ​​vụ sát hại người chủ tỷ phú của mình và bị gán tội oan, Cole Reed lên một con tàu chở hàng để đơn độc thực hiện hành trình báo thù, nhưng rồi lại phát hiện ra một âm mưu mang tầm quốc tế.",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/nDcgl5K2MCZk3xK3bh383CmIZsS.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/e2QAGrEmbpmZpMymDRkDisJkvg9.jpg",
    release_year: 2026,
    rating: 7.3,
    vote_count: 699,
    runtime_minutes: 95,
    genres: ["Phim Hành Động", "Phim Gây Cấn"],
    directors: ["Jean-François Richet"],
    cast: [],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=2Iqvbe98Gb4",
        youtube_key: "2Iqvbe98Gb4",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/e2QAGrEmbpmZpMymDRkDisJkvg9.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/hpBGCnzOvdtQoMyE48gvwp2y5yx.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/4D1pdB27uph7J8HQzNf8QvvH9bn.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1375646,
    title: "Bầy Xác Sống",
    original_title: "군체",
    overview:
      "Giáo sư Se Jeong tham dự một hội nghị công nghệ sinh học, nhưng lại chứng kiến ​​nó biến thành thảm họa khi một loại virus đột biến nhanh chóng được giải phóng. Khi dịch bệnh lan rộng và những người nhiễm bệnh bắt đầu biến đổi, chính quyền đã phong tỏa toàn bộ cơ sở.",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/lzI70txBtLH5kCUixYMJQffKkL9.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/hpBGCnzOvdtQoMyE48gvwp2y5yx.jpg",
    release_year: 2026,
    rating: 8.1,
    vote_count: 948,
    runtime_minutes: 123,
    genres: ["Phim Kinh Dị", "Phim Hành Động", "Phim Khoa Học Viễn Tưởng"],
    directors: ["연상호"],
    cast: [
      {
        name: "전지현",
        profile_url:
          "https://image.tmdb.org/t/p/w185/qejOQBdIzN18e69yRcsiD0JQi4c.jpg",
      },
      {
        name: "구교환",
        profile_url:
          "https://image.tmdb.org/t/p/w185/aFwCpSKPGFqoiTStY6ljXnD5CwH.jpg",
      },
      {
        name: "지창욱",
        profile_url:
          "https://image.tmdb.org/t/p/w185/sBmHrO5Tn27Ot5hy0yAKniROmNb.jpg",
      },
      {
        name: "김신록",
        profile_url:
          "https://image.tmdb.org/t/p/w185/x4zo2mSbehnQY80PbuldAZ7ruGm.jpg",
      },
      {
        name: "신현빈",
        profile_url:
          "https://image.tmdb.org/t/p/w185/lU0yMFh5KgOEWZRThLYbB6ZI1hE.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=LW6dpj1uCK8",
        youtube_key: "LW6dpj1uCK8",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/hpBGCnzOvdtQoMyE48gvwp2y5yx.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/4D1pdB27uph7J8HQzNf8QvvH9bn.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/5lTZyuBTNOfawsfPT8Q0cIg6qAF.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1084244,
    title: "Câu Chuyện Đồ Chơi 5",
    original_title: "Toy Story 5",
    overview:
      "Khi Bonnie được tặng chiếc máy tính bảng Lilypad và dần mê mẩn món đồ chơi mới, Buzz, Woody, Jessie cùng cả nhóm bỗng đứng trước nguy cơ bị lãng quên. Để giành lại giờ chơi, họ phải đối đầu với “đối thủ” hiện đại nhất từ trước đến nay: một màn hình nhỏ nhưng có sức hút quá lớn.",
    tagline: "Đã sẵn sàng.",
    poster_url:
      "https://image.tmdb.org/t/p/w500/8ZE3s8PN5yRpTidymRWpI0jvs7j.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/4D1pdB27uph7J8HQzNf8QvvH9bn.jpg",
    release_year: 2026,
    rating: 8.3,
    vote_count: 2409,
    runtime_minutes: 102,
    genres: ["Phim Hoạt Hình", "Phim Gia Đình", "Phim Hài", "Phim Phiêu Lưu"],
    directors: ["Andrew Stanton"],
    cast: [
      {
        name: "Joan Cusack",
        profile_url:
          "https://image.tmdb.org/t/p/w185/59UIeHZFYrKyP20lXqijtfTXglO.jpg",
      },
      {
        name: "Tom Hanks",
        profile_url:
          "https://image.tmdb.org/t/p/w185/oFvZoKI6lvU03n4YoNGAll9rkas.jpg",
      },
      {
        name: "Tim Allen",
        profile_url:
          "https://image.tmdb.org/t/p/w185/k8XPNHBO459FaqkqqM4tfX5GMKk.jpg",
      },
      {
        name: "Conan O'Brien",
        profile_url:
          "https://image.tmdb.org/t/p/w185/deRbViPut0t80miscBpP2DhBJU5.jpg",
      },
      {
        name: "Scarlett Spears",
        profile_url:
          "https://image.tmdb.org/t/p/w185/h0loJ4v3wbUfUA6dxUYv7B6fET1.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=QftAW9TTmuQ",
        youtube_key: "QftAW9TTmuQ",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/4D1pdB27uph7J8HQzNf8QvvH9bn.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/5lTZyuBTNOfawsfPT8Q0cIg6qAF.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/c6BPbkO5Npt1OdwttAxCFo06wtH.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1339713,
    title: "Ám Ảnh",
    original_title: "Obsession",
    overview:
      'Sau khi bẻ gãy "Liễu Ước Nguyện" thần bí để có được người mình thầm yêu, gã si tình rốt cuộc cũng cầu được ước thấy, để rồi kinh hãi nhận ra cái giá tăm tối phía sau lời ước đó.',
    tagline: "Cẩn thận với kẻ bạn đang cuồng si...",
    poster_url:
      "https://image.tmdb.org/t/p/w500/15qLAM3QM8DoPL9Fps4JOTdqWqt.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/5lTZyuBTNOfawsfPT8Q0cIg6qAF.jpg",
    release_year: 2026,
    rating: 8.2,
    vote_count: 5809,
    runtime_minutes: 109,
    genres: ["Phim Kinh Dị", "Phim Gây Cấn"],
    directors: ["Curry Barker"],
    cast: [
      {
        name: "Michael Johnston",
        profile_url:
          "https://image.tmdb.org/t/p/w185/fbpcCkBzu43kMdlXxEAMuLhseL8.jpg",
      },
      {
        name: "Inde Navarrette",
        profile_url:
          "https://image.tmdb.org/t/p/w185/ayKQBpXkjtLWIKOaITO493NTqFk.jpg",
      },
      {
        name: "Cooper Tomlinson",
        profile_url:
          "https://image.tmdb.org/t/p/w185/vBMQbYT1DyWPCUp11dIiqZR9zhd.jpg",
      },
      {
        name: "Megan Lawless",
        profile_url:
          "https://image.tmdb.org/t/p/w185/6qW63YEgB1qro01sM7T2HvhtFkh.jpg",
      },
      {
        name: "Andy Richter",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5Qr7N6TzC8cI0ULxDm6EC5GpZ4C.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=gMC8kkwbIQQ",
        youtube_key: "gMC8kkwbIQQ",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/5lTZyuBTNOfawsfPT8Q0cIg6qAF.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/c6BPbkO5Npt1OdwttAxCFo06wtH.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/7RyHsO4yDXtBv1zUU3mTpHeQ0d5.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 119,
        name: "Amazon Prime Video",
        logo_url:
          "https://image.tmdb.org/t/p/w92/gMZdpavHmxFNnLpMHwVxfqeux2g.png",
        type: "flatrate",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1108427,
    title: "Hành Trình Của Moana",
    original_title: "Moana",
    overview:
      "Nghe theo tiếng gọi của Đại Dương, Moana dũng cảm rời đảo Motunui, lần đầu vượt qua rạn san hô để khám phá thế giới rộng lớn. Cùng á thần Maui, cô bắt đầu hành trình phiêu lưu đầy kỳ diệu với sứ mệnh mang lại sự thịnh vượng cho quê hương.",
    tagline: "Một làn sóng mới đang trỗi dậy.",
    poster_url:
      "https://image.tmdb.org/t/p/w500/lwqKqiBsGVBtQBZF9ZTunthBCyI.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/c6BPbkO5Npt1OdwttAxCFo06wtH.jpg",
    release_year: 2026,
    rating: 7.4,
    vote_count: 805,
    runtime_minutes: 115,
    genres: ["Phim Gia Đình", "Phim Giả Tượng", "Phim Hài", "Phim Phiêu Lưu"],
    directors: ["Thomas Kail"],
    cast: [
      {
        name: "Catherine Lagaʻaia",
        profile_url:
          "https://image.tmdb.org/t/p/w185/2KRIRDwy1CtY7Bge3aqVZrORelc.jpg",
      },
      {
        name: "Dwayne Johnson",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5QApZVV8FUFlVxQpIK3Ew6cqotq.jpg",
      },
      {
        name: "Rena Owen",
        profile_url:
          "https://image.tmdb.org/t/p/w185/648ZdDBmlx6OFDFRmgAbh6q5LBo.jpg",
      },
      {
        name: "John Tui",
        profile_url:
          "https://image.tmdb.org/t/p/w185/2jIc9M5kl2GmK8fZtbtUr2s1jkS.jpg",
      },
      {
        name: "Frankie Adams",
        profile_url:
          "https://image.tmdb.org/t/p/w185/aAUHUSf0lh3OBRoaiCRL9ep8lfL.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=EEz5xbzYPKI",
        youtube_key: "EEz5xbzYPKI",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/c6BPbkO5Npt1OdwttAxCFo06wtH.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/7RyHsO4yDXtBv1zUU3mTpHeQ0d5.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/kkcwhgSFd81QDlXo8ytrpHPQjhy.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
    ],
  },
  {
    id: 299534,
    title: "Avengers: Hồi Kết",
    original_title: "Avengers: Endgame",
    overview:
      "Sau sự kiện hủy diệt tàn khốc, vũ trụ chìm trong cảnh hoang tàn. Với sự trợ giúp của những đồng minh còn sống sót, biệt đội siêu anh hùng Avengers tập hợp một lần nữa để đảo ngược hành động của Thanos và khôi phục lại trật tự của vũ trụ.",
    tagline: "Báo thù cho những người đã hi sinh.",
    poster_url:
      "https://image.tmdb.org/t/p/w500/8go3YE9sBMQaCXEx23j6BAfeuxd.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/7RyHsO4yDXtBv1zUU3mTpHeQ0d5.jpg",
    release_year: 2019,
    rating: 8.2,
    vote_count: 28841,
    runtime_minutes: 180,
    genres: ["Phim Phiêu Lưu", "Phim Khoa Học Viễn Tưởng", "Phim Hành Động"],
    directors: ["Anthony Russo", "Joe Russo"],
    cast: [
      {
        name: "Robert Downey Jr.",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5qHNjhtjMD4YWH3UP0rm4tKwxCL.jpg",
      },
      {
        name: "Chris Evans",
        profile_url:
          "https://image.tmdb.org/t/p/w185/3bOGNsHlrswhyW79uvIHH1V43JI.jpg",
      },
      {
        name: "Mark Ruffalo",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5GilHMOt5PAQh6rlUKZzGmaKEI7.jpg",
      },
      {
        name: "Chris Hemsworth",
        profile_url:
          "https://image.tmdb.org/t/p/w185/piQGdoIQOF3C1EI5cbYZLAW1gfj.jpg",
      },
      {
        name: "Scarlett Johansson",
        profile_url:
          "https://image.tmdb.org/t/p/w185/tgxYh3jMs5bY2Ub4d2dcp9iaz1R.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=4sZj4aeYUCA",
        youtube_key: "4sZj4aeYUCA",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/7RyHsO4yDXtBv1zUU3mTpHeQ0d5.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/kkcwhgSFd81QDlXo8ytrpHPQjhy.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/dqmMWNWfLnExDRpMtIMqI97GQFR.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1315772,
    title: "Minions & Quái Vật",
    original_title: "Minions & Monsters",
    overview:
      "Đây là câu chuyện có thật, đầy náo nhiệt, lố bịch và vô cùng điên rồ về cách các Minion chinh phục Hollywood, trở thành ngôi sao điện ảnh, rồi đánh mất tất cả, vô tình giải phóng những quái vật ra thế giới và cuối cùng buộc phải sát cánh bên nhau để cứu lấy Trái Đất khỏi chính mớ hỗn loạn mà họ đã gây ra.",
    tagline: "",
    poster_url:
      "https://image.tmdb.org/t/p/w500/z8OWDTR7pQuZi7jkEuR7yMXRrQt.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/kkcwhgSFd81QDlXo8ytrpHPQjhy.jpg",
    release_year: 2026,
    rating: 7.6,
    vote_count: 1192,
    runtime_minutes: 90,
    genres: [
      "Phim Phiêu Lưu",
      "Phim Hoạt Hình",
      "Phim Hài",
      "Phim Gia Đình",
      "Phim Giả Tượng",
    ],
    directors: ["Pierre Coffin"],
    cast: [
      {
        name: "Pierre Coffin",
        profile_url:
          "https://image.tmdb.org/t/p/w185/eAA9uWRqHlm1LT3nZfXb7UuPfVb.jpg",
      },
      {
        name: "Trey Parker",
        profile_url:
          "https://image.tmdb.org/t/p/w185/3tJe8t90lrcQzXj9cRcpW9wmj2c.jpg",
      },
      {
        name: "Christoph Waltz",
        profile_url:
          "https://image.tmdb.org/t/p/w185/jMvLGCVXLaBqjRLf5olyvEucZob.jpg",
      },
      {
        name: "Allison Janney",
        profile_url:
          "https://image.tmdb.org/t/p/w185/kSaSnQ9xU8eVNL0mTppab95dZA8.jpg",
      },
      {
        name: "Jesse Eisenberg",
        profile_url:
          "https://image.tmdb.org/t/p/w185/2ojZrkt5rdkUi857WSeCatxXdGS.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=V-O-uBaHk3c",
        youtube_key: "V-O-uBaHk3c",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/kkcwhgSFd81QDlXo8ytrpHPQjhy.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/dqmMWNWfLnExDRpMtIMqI97GQFR.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/mDfJG3LC3Dqb67AZ52x3Z0jU0uB.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
  {
    id: 1083381,
    title: "Backrooms: Thực Thể Quỷ Quyệt",
    original_title: "Backrooms",
    overview:
      "Backrooms theo chân Clark, một chủ cửa hàng nội thất, vô tình phát hiện cánh cửa bí ẩn dưới tầng hầm. Bước qua đó, anh bị cuốn vào một chiều không gian vô tận với những căn phòng màu vàng méo mó, liên tục lặp lại. Khi Clark ngày càng lún sâu và ám ảnh, nhà trị liệu tâm lý của anh là Mary quyết định bước vào không gian đó để tìm và giải cứu anh.",
    tagline: "Lẽ ra không nên có mặt ở đây.",
    poster_url:
      "https://image.tmdb.org/t/p/w500/7uabN7xF4209cDUdBuMvkVcpkN2.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/dqmMWNWfLnExDRpMtIMqI97GQFR.jpg",
    release_year: 2026,
    rating: 7,
    vote_count: 3487,
    runtime_minutes: 111,
    genres: ["Phim Kinh Dị", "Phim Bí Ẩn", "Phim Khoa Học Viễn Tưởng"],
    directors: ["Kane Parsons"],
    cast: [
      {
        name: "Chiwetel Ejiofor",
        profile_url:
          "https://image.tmdb.org/t/p/w185/kq5DDnqqofoRI0t6ddtRlsJnNPT.jpg",
      },
      {
        name: "Renate Reinsve",
        profile_url:
          "https://image.tmdb.org/t/p/w185/q0ljTd4fyFSJKxPgmcrKtdg3HKo.jpg",
      },
      {
        name: "Finn Bennett",
        profile_url:
          "https://image.tmdb.org/t/p/w185/p4ya77nmLlhS2cyPKUZi9zVD4Mu.jpg",
      },
      {
        name: "Lukita Maxwell",
        profile_url:
          "https://image.tmdb.org/t/p/w185/g5g6XxUAtSdPfsJrSoFkl7vH51d.jpg",
      },
      {
        name: "Mark Duplass",
        profile_url:
          "https://image.tmdb.org/t/p/w185/lRDf99rAfcdqt8Cqk4LsIT7XSD2.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=0HjdiohVOik",
        youtube_key: "0HjdiohVOik",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/dqmMWNWfLnExDRpMtIMqI97GQFR.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/mDfJG3LC3Dqb67AZ52x3Z0jU0uB.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/qeQJx07rK2xm8SD2sJxFKhE7gs0.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 119,
        name: "Amazon Prime Video",
        logo_url:
          "https://image.tmdb.org/t/p/w92/gMZdpavHmxFNnLpMHwVxfqeux2g.png",
        type: "flatrate",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
  {
    id: 299536,
    title: "Avengers: Cuộc Chiến Vô Cực",
    original_title: "Avengers: Infinity War",
    overview:
      "Sau chuyến hành trình độc nhất vô nhị không ngừng mở rộng và phát triển vụ trũ điện ảnh Marvel, bộ phim Avengers: Cuộc Chiến Vô Cực sẽ mang đến màn ảnh trận chiến cuối cùng khốc liệt nhất mọi thời đại. Biệt đội Avengers và các đồng minh siêu anh hùng của họ phải chấp nhận hy sinh tất cả để có thể chống lại kẻ thù hùng mạnh Thanos trước tham vọng hủy diệt toàn bộ vũ trụ của hắn.",
    tagline: "Cả một vũ trụ. Một lần và mãi mãi.",
    poster_url:
      "https://image.tmdb.org/t/p/w500/8gHc1cthgTOkmMiOREodCVZgJ7P.jpg",
    backdrop_url:
      "https://image.tmdb.org/t/p/w1280/mDfJG3LC3Dqb67AZ52x3Z0jU0uB.jpg",
    release_year: 2018,
    rating: 8.2,
    vote_count: 33012,
    runtime_minutes: 149,
    genres: ["Phim Phiêu Lưu", "Phim Hành Động", "Phim Khoa Học Viễn Tưởng"],
    directors: ["Anthony Russo", "Joe Russo"],
    cast: [
      {
        name: "Robert Downey Jr.",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5qHNjhtjMD4YWH3UP0rm4tKwxCL.jpg",
      },
      {
        name: "Chris Evans",
        profile_url:
          "https://image.tmdb.org/t/p/w185/3bOGNsHlrswhyW79uvIHH1V43JI.jpg",
      },
      {
        name: "Chris Hemsworth",
        profile_url:
          "https://image.tmdb.org/t/p/w185/piQGdoIQOF3C1EI5cbYZLAW1gfj.jpg",
      },
      {
        name: "Josh Brolin",
        profile_url:
          "https://image.tmdb.org/t/p/w185/sX2etBbIkxRaCsATyw5ZpOVMPTD.jpg",
      },
      {
        name: "Mark Ruffalo",
        profile_url:
          "https://image.tmdb.org/t/p/w185/5GilHMOt5PAQh6rlUKZzGmaKEI7.jpg",
      },
    ],
    media: [
      {
        type: "video",
        url: "https://www.youtube.com/watch?v=DKqu9qc-5f4",
        youtube_key: "DKqu9qc-5f4",
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/mDfJG3LC3Dqb67AZ52x3Z0jU0uB.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/qeQJx07rK2xm8SD2sJxFKhE7gs0.jpg",
        youtube_key: null,
      },
      {
        type: "image",
        url: "https://image.tmdb.org/t/p/w780/3icyRAqgakNcQn6aDVz9libFmBA.jpg",
        youtube_key: null,
      },
    ],
    providers: [
      {
        id: 2,
        name: "Apple TV Store",
        logo_url:
          "https://image.tmdb.org/t/p/w92/qdEGArH3lKfFnAtYXMkSYk5wxuG.png",
        type: "rent",
      },
      {
        id: 192,
        name: "YouTube",
        logo_url:
          "https://image.tmdb.org/t/p/w92/5Maob4o5w8oZnNeYpCDyVFD3M7X.png",
        type: "rent",
      },
    ],
  },
];
