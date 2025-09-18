const postList = document.getElementById("post-list");
const loadMore = document.getElementById("load-more");
loadMore.style.display = "block";

let beforeId = null;
let loading = false;
let finished = false;

let timeout = null;
let throttlingTimeInterval = 200;
let loadmorePixelThreshold = 200;

async function fetchPosts() {
    if (loading || finished) return;
    loading = true;
    loadMore.textContent = "Loading...";

    try {
        const url = beforeId
            ? `/api/guest/posts?before_id=${beforeId}`
            : "/api/guest/posts";

        const response = await fetch(url, { method: "GET" });
        if (!response.ok) throw new Error("Failed to fetch posts");

        const result = await response.json();
        if (!result.success) throw new Error(result.message || "Unknown error");

        const posts = result.data.post_briefs;
        beforeId = result.data.next_before_id;

        if (posts.length === 0) {
            finished = true;
            loadMore.textContent = "All caught up";
            return;
        }

        const html = posts.map(post => {
            const category = post.category.trim() !== "" ?
                `<div class="category">
                    Category: <a href="/categories?category=${encodeURIComponent(post.category)}">${post.category}</a>
                </div>`
                : "";

            const tags = Array.isArray(post.tags) && post.tags.length > 0 ?
                `<div class="tags">Tags: ${post.tags.map(tag =>
                    `<a href="/tags?tag=${encodeURIComponent(tag)}">#${tag}</a>`
                ).join(" ")}</div>`
                : "";

            return `
                    <article class="post-item">
                        <h2><a href="/post/${encodeURIComponent(post.slug)}">${post.title}</a></h2>
                        <div class="meta">
                            ${category}
                            ${tags}
                        </div>
                    </article>
                `;
        }).join("");

        postList.insertAdjacentHTML("beforeend", html);

        if (!beforeId) {
            finished = true;
            loadMore.textContent = "All caught up";
        } else {
            loadMore.textContent = "Pull to load";
        }
    } catch (err) {
        console.error(err);
        loadMore.textContent = `Error: ${err.message}`;
    } finally {
        loading = false;
    }
}

async function fillScreenPosts() {
    if (!finished) {
        await fetchPosts();
    }
    while (!finished && document.body.offsetHeight <= window.innerHeight) {
        await fetchPosts();
    }
}

document.addEventListener("DOMContentLoaded", () => {
    fillScreenPosts();

    window.addEventListener("scroll", () => {
        if (finished || loading) return;

        if (timeout) return;
        timeout = setTimeout(() => {
            if (window.innerHeight + window.scrollY >= document.body.offsetHeight - loadmorePixelThreshold) {
                fetchPosts();
            }
        }, throttlingTimeInterval);
    });
});