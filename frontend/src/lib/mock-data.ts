// Mock API data, used when NEXT_PUBLIC_API_BASE_URL is unset (frontend/CLAUDE.md).
// Movies are copied from a real refresh of the backend cache (option "friends").
// Some fields are altered on purpose so every hide rule and fallback in
// docs/SCREENS.md shows up in mock mode:
// - 969681: has providers (logos null → name-only fallback)
// - 1368337: no poster, no backdrop
// - 1204680: no runtime, no release year
// - 1288445: empty cast
// Titles with an empty overview are real TMDB gaps (no vi-VN translation).
import type { MovieDetail, ViewingOption } from "./api-types";

export const MOCK_FETCHED_AT = "2026-09-27T17:22:13Z";

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
    vote_count: 2895,
    runtime_minutes: 150,
    genres: ["Phim Khoa Học Viễn Tưởng", "Phim Hành Động", "Phim Phiêu Lưu"],
    directors: ["Destin Daniel Cretton"],
    cast: [
      "Tom Holland",
      "Zendaya",
      "Mark Ruffalo",
      "Jon Bernthal",
      "Jacob Batalon",
    ],
    trailer_url: "https://www.youtube.com/watch?v=P3uI5sLosKU",
    providers: [
      {
        id: 8,
        name: "Netflix",
        logo_url: null,
        type: "flatrate",
      },
      {
        id: 2,
        name: "Apple TV",
        logo_url: null,
        type: "rent",
      },
    ],
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
    vote_count: 520,
    runtime_minutes: 95,
    genres: ["Phim Kinh Dị", "Phim Khoa Học Viễn Tưởng", "Phim Phiêu Lưu"],
    directors: ["Zach Cregger"],
    cast: [
      "Austin Abrams",
      "Paul Walter Hauser",
      "Kali Reis",
      "Zach Cherry",
      "Johnno Wilson",
    ],
    trailer_url: "https://www.youtube.com/watch?v=mNd1gb19A-c",
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
    vote_count: 3921,
    runtime_minutes: 173,
    genres: ["Phim Phiêu Lưu", "Phim Hành Động", "Phim Giả Tượng"],
    directors: ["Christopher Nolan"],
    cast: [
      "Matt Damon",
      "Tom Holland",
      "Anne Hathaway",
      "Robert Pattinson",
      "Himesh Patel",
    ],
    trailer_url: "https://www.youtube.com/watch?v=o4s236zGkMM",
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
    vote_count: 534,
    runtime_minutes: null,
    genres: ["Phim Hài", "Phim Phiêu Lưu", "Phim Gia Đình"],
    directors: ["Dave Green"],
    cast: [
      "Will Forte",
      "Trần Đồng Lan",
      "John Cena",
      "Tone Bell",
      "Martha Kelly",
    ],
    trailer_url: "https://www.youtube.com/watch?v=kMsiD1Nky5I",
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
    rating: 8.3,
    vote_count: 279,
    runtime_minutes: 112,
    genres: ["Phim Lãng Mạn", "Phim Hài"],
    directors: ["Claire Scanlon"],
    cast: [
      "Lili Reinhart",
      "Tom Bateman",
      "Rachel Marsh",
      "Jaboukie Young-White",
      "Nicholas Duvernay",
    ],
    trailer_url: "https://www.youtube.com/watch?v=xwdamnfS6IM",
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
    rating: 7.2,
    vote_count: 608,
    runtime_minutes: 95,
    genres: ["Phim Hành Động", "Phim Gây Cấn"],
    directors: ["Jean-François Richet"],
    cast: [],
    trailer_url: "https://www.youtube.com/watch?v=2Iqvbe98Gb4",
    providers: [],
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
    vote_count: 889,
    runtime_minutes: 123,
    genres: ["Phim Hành Động", "Phim Kinh Dị", "Phim Khoa Học Viễn Tưởng"],
    directors: ["연상호"],
    cast: ["전지현", "구교환", "지창욱", "김신록", "신현빈"],
    trailer_url: "https://www.youtube.com/watch?v=LW6dpj1uCK8",
    providers: [],
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
    vote_count: 2349,
    runtime_minutes: 102,
    genres: ["Phim Hoạt Hình", "Phim Gia Đình", "Phim Hài", "Phim Phiêu Lưu"],
    directors: ["Andrew Stanton"],
    cast: [
      "Joan Cusack",
      "Tom Hanks",
      "Tim Allen",
      "Conan O'Brien",
      "Scarlett Spears",
    ],
    trailer_url: "https://www.youtube.com/watch?v=QftAW9TTmuQ",
    providers: [],
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
      "https://image.tmdb.org/t/p/w1280/uoTZOIXnkcTgJe1emITpsQq7sPu.jpg",
    release_year: 2026,
    rating: 8.2,
    vote_count: 5751,
    runtime_minutes: 109,
    genres: ["Phim Kinh Dị", "Phim Gây Cấn"],
    directors: ["Curry Barker"],
    cast: [
      "Michael Johnston",
      "Inde Navarrette",
      "Cooper Tomlinson",
      "Megan Lawless",
      "Andy Richter",
    ],
    trailer_url: "https://www.youtube.com/watch?v=gMC8kkwbIQQ",
    providers: [],
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
    vote_count: 781,
    runtime_minutes: 115,
    genres: ["Phim Gia Đình", "Phim Giả Tượng", "Phim Hài", "Phim Phiêu Lưu"],
    directors: ["Thomas Kail"],
    cast: [
      "Catherine Lagaʻaia",
      "Dwayne Johnson",
      "Rena Owen",
      "John Tui",
      "Frankie Adams",
    ],
    trailer_url: "https://www.youtube.com/watch?v=EEz5xbzYPKI",
    providers: [],
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
    vote_count: 28809,
    runtime_minutes: 180,
    genres: ["Phim Phiêu Lưu", "Phim Khoa Học Viễn Tưởng", "Phim Hành Động"],
    directors: ["Anthony Russo", "Joe Russo"],
    cast: [
      "Robert Downey Jr.",
      "Chris Evans",
      "Mark Ruffalo",
      "Chris Hemsworth",
      "Scarlett Johansson",
    ],
    trailer_url: "https://www.youtube.com/watch?v=4sZj4aeYUCA",
    providers: [],
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
    vote_count: 1180,
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
      "Pierre Coffin",
      "Trey Parker",
      "Christoph Waltz",
      "Allison Janney",
      "Jesse Eisenberg",
    ],
    trailer_url: "https://www.youtube.com/watch?v=V-O-uBaHk3c",
    providers: [],
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
    vote_count: 3440,
    runtime_minutes: 111,
    genres: ["Phim Kinh Dị", "Phim Bí Ẩn", "Phim Khoa Học Viễn Tưởng"],
    directors: ["Kane Parsons"],
    cast: [
      "Chiwetel Ejiofor",
      "Renate Reinsve",
      "Finn Bennett",
      "Lukita Maxwell",
      "Mark Duplass",
    ],
    trailer_url: "https://www.youtube.com/watch?v=0HjdiohVOik",
    providers: [],
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
    vote_count: 32993,
    runtime_minutes: 149,
    genres: ["Phim Phiêu Lưu", "Phim Hành Động", "Phim Khoa Học Viễn Tưởng"],
    directors: ["Joe Russo", "Anthony Russo"],
    cast: [
      "Robert Downey Jr.",
      "Chris Evans",
      "Chris Hemsworth",
      "Josh Brolin",
      "Mark Ruffalo",
    ],
    trailer_url: "https://www.youtube.com/watch?v=DKqu9qc-5f4",
    providers: [],
  },
];
