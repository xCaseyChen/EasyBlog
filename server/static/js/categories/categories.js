const postList = document.getElementById("post-list");
const loadMore = document.getElementById("load-more");
const categoriesContainer = document.getElementById("categories-container");
loadMore.style.display = "block";

let beforeId = null;
let loading = false;
let finished = false;

let timeout = null;
let throttlingTimeInterval = 200;
let loadmorePixelThreshold = 200;

let curCategory = "";

function resetPostList() {
    beforeId = null;
    loading = false;
    finished = false;
    curCategory = "";
    postList.innerHTML = "";
    loadMore.textContent = "";
}

function generateCategoryId(category) {
    return "category-" + category.replace(/\s+/g, '-').replace(/[^a-zA-Z0-9-_]/g, '');
}

async function fetchCategories() {
    try {
        const url = "/api/guest/categories";

        const response = await fetch(url, { method: "GET" });
        if (!response.ok) throw new Error("Failed to fetch categories");

        const result = await response.json();
        if (!result.success) throw new Error(result.message || "Unknown error");

        const categories = result.data.categories;

        if (categories.length === 0) {
            console.log("No categories available");
            return;
        }

        categoriesContainer.innerHTML = "";

        categories.forEach(category => {
            const label = document.createElement("label");
            label.classList.add("category-label");

            const input = document.createElement("input");
            input.type = "checkbox";
            input.name = "category";
            input.value = category;
            input.id = generateCategoryId(category);

            input.addEventListener("change", (e) => {
                resetPostList();
                if (e.target.checked) {
                    const allCheckboxes = document.querySelectorAll('input[name="category"]');
                    allCheckboxes.forEach(checkbox => {
                        if (checkbox !== e.target) {
                            checkbox.checked = false;
                        }
                    });
                    curCategory = e.target.value;
                }
                fillScreenPostsByCurCategory();
            });

            label.appendChild(input);
            label.appendChild(document.createTextNode(category));

            categoriesContainer.appendChild(label);
        });
    } catch (err) {
        console.error(err);
        loadMore.textContent = `Error: ${err.message}`;
    }
}

async function fillScreenPostsByCurCategory() {
    if (curCategory !== "" && !finished) {
        await fetchPostsByCurCategory();
    }
    while (curCategory !== "" && !finished && document.body.offsetHeight <= window.innerHeight) {
        await fetchPostsByCurCategory();
    }
}

async function fetchPostsByCurCategory() {
    if (loading || finished || curCategory === "") return;
    loading = true;
    loadMore.textContent = "Loading...";

    try {
        const url = beforeId
            ? `/api/guest/posts?category=${curCategory}&before_id=${beforeId}`
            : `/api/guest/posts?category=${curCategory}`;

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
                        <h2><a href="/posts/${encodeURIComponent(post.slug)}">${post.title}</a></h2>
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

function clickCategoryFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    const category = urlParams.get('category');

    if (category) {
        const categoryCheckbox = document.querySelector(`#${generateCategoryId(category)}`);
        if (categoryCheckbox) {
            categoryCheckbox.click();
        }
    }
}

document.addEventListener("DOMContentLoaded", () => {
    fetchCategories().then(() => {
        clickCategoryFromUrl();
    });

    window.addEventListener("scroll", () => {
        if (finished || loading || curCategory === "") return;

        if (timeout) return;
        timeout = setTimeout(() => {
            if (window.innerHeight + window.scrollY >= document.body.offsetHeight - loadmorePixelThreshold) {
                fetchPostsByCurCategory();
            }
            timeout = null;
        }, throttlingTimeInterval);
    });
})