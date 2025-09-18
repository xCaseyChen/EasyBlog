const postList = document.getElementById("post-list");
const loadMore = document.getElementById("load-more");
const tagsContainer = document.getElementById("tags-container");
loadMore.style.display = "block";

let beforeId = null;
let loading = false;
let finished = false;

let timeout = null;
let throttlingTimeInterval = 200;
let loadmorePixelThreshold = 200;

let curTags = new Set();

function resetPostList() {
    beforeId = null;
    loading = false;
    finished = false;
    postList.innerHTML = "";
    loadMore.textContent = "";
}

function generateTagId(tag) {
    return "tag-" + tag.replace(/\s+/g, '-').replace(/[^a-zA-Z0-9-_]/g, '');
}

async function fetchTags() {
    try {
        const url = "/api/guest/tags";

        const response = await fetch(url, { method: "GET" });
        if (!response.ok) throw new Error("Failed to fetch tags");

        const result = await response.json();
        if (!result.success) throw new Error(result.message || "Unknown error");

        const tags = result.data.tags;

        if (tags.length === 0) {
            console.log("No tags available");
            return;
        }

        tags.forEach(tag => {
            const label = document.createElement("label");
            label.classList.add("tags-label");

            const input = document.createElement("input");
            input.type = "checkbox";
            input.name = "tag";
            input.value = tag;
            input.id = generateTagId(tag);

            input.addEventListener("change", (e) => {
                resetPostList();
                if (e.target.checked) {
                    curTags.add(tag);
                } else {
                    curTags.delete(tag)
                }
                fillScreenPostsByCurTags();
            });

            label.appendChild(input);
            label.appendChild(document.createTextNode(`#${tag}`));

            tagsContainer.appendChild(label);
        });
    } catch (err) {
        console.log(err);
        loadMore.textContent = `Error: ${err.message}`;
    }
}

async function fillScreenPostsByCurTags() {
    if (curTags.size !== 0 && !finished) {
        await fetchPostsByCurTags();
    }
    while (curTags.size !== 0 && !finished && document.body.offsetHeight <= window.innerHeight) {
        await fetchPostsByCurTags();
    }
}

async function fetchPostsByCurTags() {
    if (loading || finished || curTags.size == 0) return;
    loading = true;
    loadMore.textContent = "Loading...";

    const curTagsStr = [...curTags].join(",")

    try {
        const url = beforeId
            ? `/api/guest/posts?tags=${curTagsStr}&before_id=${beforeId}`
            : `/api/guest/posts?tags=${curTagsStr}`;

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
        finished = true;
    } finally {
        loading = false;
    }
}

function clickTagFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    const tag = urlParams.get('tag');

    if (tag) {
        const tagCheckbox = document.querySelector(`#${generateTagId(tag)}`);
        if (tagCheckbox) {
            tagCheckbox.click();
        }
    }
}

document.addEventListener("DOMContentLoaded", () => {
    fetchTags().then(() => {
        clickTagFromUrl();
    });

    window.addEventListener("scroll", () => {
        if (finished || loading || curTags.size === 0) return;

        if (timeout) return;
        timeout = setTimeout(() => {
            if (window.innerHeight + window.scrollY >= document.body.offsetHeight - loadmorePixelThreshold) {
                fetchPostsByCurTags();
            }
            timeout = null;
        }, throttlingTimeInterval);
    });
})