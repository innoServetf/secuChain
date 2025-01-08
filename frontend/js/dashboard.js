document.addEventListener('DOMContentLoaded', function() {
    // 检查登录状态
    if (!localStorage.getItem('token')) {
        window.location.href = '../index.html';
        return;
    }

    // 加载用户信息
    loadUserInfo();
    // 加载统计数据
    loadDashboardStats();
    // 加载最近的SBOM
    loadRecentSboms();

    // 退出登录
    document.getElementById('logoutBtn').addEventListener('click', logout);
});

async function loadUserInfo() {
    try {
        const response = await Api.getUserInfo();
        const user = response.data;
        
        document.getElementById('userInfo').textContent = user.username;
        document.getElementById('organizationInfo').textContent = 
            `组织: ${user.organization}`;
    } catch (error) {
        console.error('Failed to load user info:', error);
        if (error.message === 'unauthorized') {
            logout();
        }
    }
}

async function loadDashboardStats() {
    try {
        const response = await Api.request('/dashboard/stats');
        const stats = response.data;

        document.getElementById('sbomCount').textContent = stats.sbomCount;
        document.getElementById('vulnerabilityCount').textContent = stats.vulnerabilityCount;
        document.getElementById('componentCount').textContent = stats.componentCount;
    } catch (error) {
        console.error('Failed to load dashboard stats:', error);
    }
}

async function loadRecentSboms() {
    try {
        const response = await Api.request('/sbom/recent');
        const sboms = response.data;

        const tbody = document.querySelector('#recentSbomTable tbody');
        tbody.innerHTML = '';

        sboms.forEach(sbom => {
            const tr = document.createElement('tr');
            tr.innerHTML = `
                <td>${sbom.name}</td>
                <td>${sbom.version}</td>
                <td>${formatDate(sbom.createdAt)}</td>
                <td><span class="status-badge ${sbom.status}">${sbom.status}</span></td>
                <td>
                    <button class="btn btn-text" onclick="viewSbom(${sbom.id})">查看</button>
                </td>
            `;
            tbody.appendChild(tr);
        });
    } catch (error) {
        console.error('Failed to load recent SBOMs:', error);
    }
}

function formatDate(dateString) {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    });
}

function viewSbom(id) {
    window.location.href = `sbom-detail.html?id=${id}`;
}

function logout() {
    localStorage.removeItem('token');
    window.location.href = '../index.html';
} 